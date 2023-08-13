#ifndef LIB_CALCMODELWRAPPER_H
#define LIB_CALCMODELWRAPPER_H

#ifdef __cplusplus
extern "C" {
#endif

struct Result {
  double value;
  const char* error;
};

struct Result calculate(const char* str, double x);
#ifdef __cplusplus
}
#endif

#endif  // LIB_CALCMODELWRAPPER_H
