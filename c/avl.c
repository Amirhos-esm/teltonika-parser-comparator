#include "avl.h"

// CRC-16/IBM (XMODEM) polynomial: x^16 + x^12 + x^5 + 1

uint16_t calculate_crc16(const uint8_t *data, size_t start, size_t end, size_t size) {
    uint16_t crc = 0x0000; // Initial CRC value
    if (data == NULL)
        return 0;
    for (size_t i = start; i != end; i++) {
        crc ^= (uint16_t) data[i];
        for (unsigned k = 0; k < 8; k++) {
            crc = crc & 1 ? (crc >> 1) ^ 0xa001 : crc >> 1;
        }
    }
    return crc;

}

bool avl_add_n1(Avl_Data_t *data, uint8_t io_id, uint8_t io_value) {
    if (data->io.n_of_n1 >= sizeof(data->io.n1_id))
        return false;
    data->io.n1_id[data->io.n_of_n1] = io_id;
    data->io.n1[data->io.n_of_n1] = io_value;
    data->io.n_of_n1++;
    return true;
}

bool avl_add_n2(Avl_Data_t *data, uint8_t io_id, uint16_t io_value) {
    if (data->io.n_of_n2 >= sizeof(data->io.n2_id))
        return false;
    data->io.n2_id[data->io.n_of_n2] = io_id;
    data->io.n2[data->io.n_of_n2] = io_value;
    data->io.n_of_n2++;
    return true;
}

bool avl_add_n4(Avl_Data_t *data, uint8_t io_id, uint32_t io_value) {
    if (data->io.n_of_n4 >= sizeof(data->io.n4_id))
        return false;
    data->io.n4_id[data->io.n_of_n4] = io_id;
    data->io.n4[data->io.n_of_n4] = io_value;
    data->io.n_of_n4++;
    return true;
}

bool avl_add_n8(Avl_Data_t *data, uint8_t io_id, uint64_t io_value) {
    if (data->io.n_of_n8 >= sizeof(data->io.n8_id))
        return false;
    data->io.n8_id[data->io.n_of_n8] = io_id;
    data->io.n8[data->io.n_of_n8] = io_value;
    data->io.n_of_n8++;
    return true;
}


bool avl_encode_avl_data(OStream *stream, Avl_Data_t *data) {
    OStream_setByteOrder(stream, ByteOrder_BigEndian);
    if (OStream_writeUInt64(stream, data->timestamp) != Stream_Ok) return false;
    if (OStream_writeUInt8(stream, (uint8_t) data->priority) != Stream_Ok) return false;

    // gps elmenets
    if (OStream_writeInt32(stream, (int32_t) (data->gps.longitude * 1e7)) != Stream_Ok) return false;
    if (OStream_writeInt32(stream, (int32_t) (data->gps.latitude * 1e7)) != Stream_Ok) return false;
    if (OStream_writeUInt16(stream, data->gps.altitude) != Stream_Ok) return false;
    if (OStream_writeUInt16(stream, data->gps.angle) != Stream_Ok) return false;
    if (OStream_writeUInt8(stream, data->gps.satellites) != Stream_Ok) return false;
    if (OStream_writeUInt16(stream, data->gps.speed) != Stream_Ok) return false;
    // io elements
    if (OStream_writeUInt8(stream, data->io.event_io_id) != Stream_Ok) return false;
    if (OStream_writeUInt8(stream, data->io.n_of_n1 + data->io.n_of_n2 + data->io.n_of_n4 + data->io.n_of_n8) !=
        Stream_Ok)
        return false;

    //n1
    if(data->io.n_of_n1 > sizeof(data->io.n1_id)){
        return false;
    }
    if (OStream_writeUInt8(stream, data->io.n_of_n1) != Stream_Ok) return false;
    for (int i = 0; i < data->io.n_of_n1; i++) {
        if (OStream_writeUInt8(stream, data->io.n1_id[i]) != Stream_Ok) return false;
        if (OStream_writeUInt8(stream, data->io.n1[i]) != Stream_Ok) return false;
    }

    //n2
    if(data->io.n_of_n2 > sizeof(data->io.n2_id)){
        return false;
    }
    if (OStream_writeUInt8(stream, data->io.n_of_n2) != Stream_Ok) return false;
    for (int i = 0; i < data->io.n_of_n2; i++) {
        if (OStream_writeUInt8(stream, data->io.n2_id[i]) != Stream_Ok) return false;
        if (OStream_writeUInt16(stream, data->io.n2[i]) != Stream_Ok) return false;
    }

    //n4
    if(data->io.n_of_n4 > sizeof(data->io.n4_id)){
        return false;
    }
    if (OStream_writeUInt8(stream, data->io.n_of_n4) != Stream_Ok) return false;
    for (int i = 0; i < data->io.n_of_n4; i++) {
        if (OStream_writeUInt8(stream, data->io.n4_id[i]) != Stream_Ok) return false;
        if (OStream_writeUInt32(stream, data->io.n4[i]) != Stream_Ok) return false;
    }


    //n8
    if(data->io.n_of_n8 > sizeof(data->io.n8_id)){
        return false;
    }
    if (OStream_writeUInt8(stream, data->io.n_of_n8) != Stream_Ok) return false;
    for (int i = 0; i < data->io.n_of_n8; i++) {
        if (OStream_writeUInt8(stream, data->io.n8_id[i]) != Stream_Ok) return false;
        if (OStream_writeUInt64(stream, data->io.n8[i]) != Stream_Ok) return false;
    }
    return true;

}


bool avl_encode(OStream *stream, Avl_Data_t *data) {

    OStream_setByteOrder(stream, ByteOrder_BigEndian);
    // preamble
    if (OStream_writeUInt32(stream, 0) != Stream_Ok) return false;
    // Data Field Length
    uint32_t data_field_length = 3 + 8 + 1 + 15;
    data_field_length +=
            1 + 1 + 1 + data->io.n_of_n1 * (1 + 1) + 1 + data->io.n_of_n2 * (1 + 2) + 1 + data->io.n_of_n4 * (1 + 4)
            + 1 + data->io.n_of_n8 * (1 + 8);
    if (OStream_writeUInt32(stream, data_field_length) != Stream_Ok) return false;

    // codec id
    size_t startPos = stream->Buffer.WPos;
    if (OStream_writeUInt8(stream, 0x08) != Stream_Ok) return false;
//    printf("start %x\n",stream->Buffer.Data[startPos]);

    // number of record
    if (OStream_writeUInt8(stream, 0x01) != Stream_Ok) return false;
    if (!avl_encode_avl_data(stream, data)) {
        return false;
    }
    // number of record #2
    if (OStream_writeUInt8(stream, 0x01) != Stream_Ok) return false;
    size_t endPos = stream->Buffer.WPos;
//    printf("diff: %d\n",endPos-startPos);


    if( OStream_writeUInt8(stream, calculate_crc16(
            stream->Buffer.Data, startPos, endPos, stream->Buffer.Size
    )) != Stream_Ok ){
      return false;
    };
    return true;
}

uint32_t getDataFeildSize(Avl_Data_t **array, uint8_t record_count) {
    uint32_t data_field_length = 3;
    for (int i = 0; i < record_count; i++) {
        Avl_Data_t *data = array[i];
        data_field_length += 8 + 1 + 15 +
                             1 + 1 + 1 + data->io.n_of_n1 * (1 + 1) + 1 + data->io.n_of_n2 * (1 + 2) + 1 +
                             data->io.n_of_n4 * (1 + 4)
                             + 1 + data->io.n_of_n8 * (1 + 8);
    }
    return data_field_length;
}

bool avl_encode_multi(OStream *stream, Avl_Data_t **data, uint8_t record_count) {
    if (record_count == 0) {
        return false;
    }
    OStream_setByteOrder(stream, ByteOrder_BigEndian);
    // preamble
    if (OStream_writeUInt32(stream, 0) != Stream_Ok) return false;
    // Data Field Length
    if (OStream_writeUInt32(stream, getDataFeildSize(data, record_count)) != Stream_Ok) return false;

    // codec id
    size_t startPos = stream->Buffer.WPos;
    if (OStream_writeUInt8(stream, 0x08) != Stream_Ok) return false;

    // number of record
    if (OStream_writeUInt8(stream, record_count) != Stream_Ok) return false;

    for (int i = 0; i < record_count; i++) {
        if (!avl_encode_avl_data(stream, data[i])) {
            return false;
        }
    }
    // number of record #2
    if (OStream_writeUInt8(stream, record_count) != Stream_Ok) return false;
    size_t endPos = stream->Buffer.WPos;


    OStream_writeUInt32(stream, calculate_crc16(
            stream->Buffer.Data, startPos, endPos, stream->Buffer.Size
    ));
    return true;
}
int avl_encode_avl_data_buffer(uint8_t *buffer, int len, Avl_Data_t *data){
    OStream stream;
    OStream_init(&stream, NULL, buffer, len);
    if (avl_encode_avl_data(&stream, data)) {
        return stream.Buffer.WPos;
    }
    return -1;
}
int avl_encode_into_buffer(uint8_t *buffer, int len, Avl_Data_t *data) {
    OStream stream;
    OStream_init(&stream, NULL, buffer, len);
    if (avl_encode(&stream, data)) {
        return stream.Buffer.WPos;
    }
    return -1;
}

int avl_encode_multi_buffer(uint8_t *buffer, int len, Avl_Data_t **data, uint8_t record_count) {
    OStream stream;
    OStream_init(&stream, NULL, buffer, len);
    if (avl_encode_multi(&stream, data, record_count)) {
        return stream.Buffer.WPos;
    }
    return -1;
}
/*
 * avl_datas show have extra space before and after it.
 *
 * */
bool avl_encode_avl_data_to_packet(uint8_t * avl_datas,int number_od_records,int len){
    /*
     * warning: accessing to outbond of avl_datas. make sure that it free and alocated
     * */
    OStream stream;
    OStream_init(&stream, NULL, avl_datas - 10, 10);
    OStream_setByteOrder(&stream, ByteOrder_BigEndian);
    // preamble
    if (OStream_writeUInt32(&stream, 0) != Stream_Ok) return false;
    if (OStream_writeUInt32(&stream, len + 3) != Stream_Ok) return false;
    if (OStream_writeUInt8(&stream, 0x08) != Stream_Ok) return false;

    if (OStream_writeUInt8(&stream, number_od_records) != Stream_Ok) return false;


    OStream_init(&stream, NULL, avl_datas + len, 5);
    OStream_setByteOrder(&stream, ByteOrder_BigEndian);
    if (OStream_writeUInt8(&stream, number_od_records) != Stream_Ok) return false;
    if( OStream_writeUInt32(&stream, calculate_crc16(
            avl_datas - 2, 0, 2 + len + 1, 0
    )) != Stream_Ok ){
        return false;
    };
    return true;

}