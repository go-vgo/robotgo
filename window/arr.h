// #include <stdio.h>
// #include <stdlib.h>
// #include <stdbool.h>

struct Arr{
    int len;
    int cnu;
    char** pName;
    int *pId;
};

// int init_Arr(struct Arr *pArray, int len){
//     pArray->pName = (char**)malloc(sizeof(char*)*len);
//     pArray->pId = (int*)malloc(sizeof(int)*len);
//     if(NULL == pArray->pName){
//         printf("false \n");
//         return 1;
//     }else{
//         pArray->len = len;
//         pArray->cnu = 0;
//         printf("true %d \n", pArray->len);
//         return 0;
//     }

//     return 0;
// }