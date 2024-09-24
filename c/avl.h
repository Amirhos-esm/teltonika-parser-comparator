#ifndef SERIALPORT_AVL_H
#define SERIALPORT_AVL_H
#include "Stream/InputStream.h"
#include "Stream/OutputStream.h"
#include "comp.h"

typedef  enum {
    Priority_low = 0,
    Priority_high =1,
    Priority_panic=2
}Priority_t;
typedef struct {
    double latitude; // in decimal degree
    double longitude; // in decimal degree
    uint16_t altitude;
    uint16_t angle;
    uint8_t satellites;
    uint16_t speed;

}Gps_Element_t;
typedef struct {
    uint8_t event_io_id;
    uint8_t n_of_n1;
    uint8_t n1_id[16];
    uint8_t n1[16];
    uint8_t n_of_n2;
    uint8_t n2_id[8];
    uint16_t n2[8];
    uint8_t n_of_n4;
    uint8_t n4_id[4];
    uint32_t n4[4];
    uint8_t n_of_n8;
    uint8_t n8_id[2];
    uint64_t n8[2];
}Io_Element_t;
typedef struct {
    uint64_t timestamp; // in ms
    Priority_t priority;
    Gps_Element_t gps;
    Io_Element_t io;
}Avl_Data_t;

bool avl_add_n1(Avl_Data_t *data,uint8_t io_id,uint8_t io_value);
bool avl_add_n2(Avl_Data_t *data,uint8_t io_id,uint16_t io_value);
bool avl_add_n4(Avl_Data_t *data,uint8_t io_id,uint32_t io_value);
bool avl_add_n8(Avl_Data_t *data,uint8_t io_id,uint64_t io_value);

/*
 * avl_datas show have extra space before and after it.
 *
 * */
bool avl_encode_avl_data_to_packet(uint8_t * avl_datas,int number_od_records,int len);
bool avl_encode_avl_data(OStream *stream, Avl_Data_t *data);
int avl_encode_avl_data_buffer(uint8_t *buffer, int len, Avl_Data_t *data);
int avl_encode_into_buffer(uint8_t *buffer, int len, Avl_Data_t *data);
bool avl_encode(OStream *stream,Avl_Data_t *data);

bool avl_encode_multi(OStream *stream,Avl_Data_t **data,uint8_t record_count);
int avl_encode_multi_buffer(uint8_t *buffer, int len, Avl_Data_t **data, uint8_t record_count);

#endif //SERIALPORT_AVL_H
