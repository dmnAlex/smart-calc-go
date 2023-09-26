#ifndef LIB_CALCMODEL_H
#define LIB_CALCMODEL_H

#include <cmath>
#include <cstring>
#include <stack>
#include <string>

enum err_t {
  ERR_SUCCESS,
  ERR_MALFORMED,
  ERR_DIVBYZERO,
  ERR_MISSBRACKET,
  ERR_NEGATIVEROOT,
};

enum val_t {
  // Precendence 0
  VAL_UNKNOWN,
  VAL_LBRACKET,
  VAL_RBRACKET,
  VAL_NUMBER,
  VAL_VARIABLE,
  // Precendence 1
  VAL_PLUS,
  VAL_MINUS,
  // Precendence 2
  VAL_MULTIPLY,
  VAL_DIVIDE,
  VAL_MOD,
  // Precendence 3
  VAL_POWER,
  // Precendence 4
  VAL_COS,
  VAL_SIN,
  VAL_TAN,
  VAL_ACOS,
  VAL_ASIN,
  VAL_ATAN,
  VAL_SQRT,
  VAL_LN,
  VAL_LOG,
};

struct item_t {
  double value;
  val_t type;
};

class CalcModel {
 public:
  CalcModel();

  double get_result();
  err_t calculate(const std::string& src, double x);

 private:
  item_t stack_pop(std::stack<item_t>& stack);
  err_t stack_push(std::stack<item_t>& stack, item_t item);
  item_t stack_peek(std::stack<item_t>& stack);
  void stack_reverse(std::stack<item_t>& stack);
  void stack_clean(std::stack<item_t>& stack);

  int get_precedence(val_t type);

  err_t parse_token(char** src, size_t& index, item_t& item);
  err_t parse_rpn(const char* src);
  err_t evaluate_rpn();
  err_t math(item_t op);

  const char* op_values[20] = {"unknown", "(",    ")",    "",    "x",
                               "+",       "-",    "*",    "/",   "mod",
                               "^",       "cos",  "sin",  "tan", "acos",
                               "asin",    "atan", "sqrt", "ln",  "log"};

  std::stack<item_t> main_stack, oper_stack;
  double x_value, out_value;
};

#endif  // LIB_CALCMODEL_H
