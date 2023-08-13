#include "calcmodelwrapper.h"

#include <iostream>

#include "calcmodel.h"

extern "C" {
struct Result calculate(const char* str, double x) {
  static CalcModel calculator;

  std::string equation(str);
  auto err = calculator.calculate(equation, x);
  Result res = {0.0, ""};

  switch (err) {
    case ERR_SUCCESS:
      res.value = calculator.get_result();
      break;
    case ERR_MALFORMED:
      res.error = "Malformed expression";
      break;
    case ERR_DIVBYZERO:
      res.error = "Division by zero";
      break;
    case ERR_MISSBRACKET:
      res.error = "Mismatched brackets";
      break;
    case ERR_NEGATIVEROOT:
      res.error = "Negative root";
      break;
    default:
      res.error = "Unknown error";
      break;
  }

  return res;
}
}
