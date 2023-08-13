#include "calcmodel.h"

CalcModel::CalcModel() : x_value{0.0}, out_value{0.0} {}

double CalcModel::get_result() { return out_value; }

item_t CalcModel::stack_pop(std::stack<item_t>& stack) {
  item_t result = {0.0, VAL_UNKNOWN};

  if (!stack.empty()) {
    result = stack.top();
    stack.pop();
  }

  return result;
}

err_t CalcModel::stack_push(std::stack<item_t>& stack, item_t item) {
  stack.push(item);

  return ERR_SUCCESS;
}

item_t CalcModel::stack_peek(std::stack<item_t>& stack) {
  item_t result = {0.0, VAL_UNKNOWN};

  if (!stack.empty()) {
    result = stack.top();
  }

  return result;
}

void CalcModel::stack_reverse(std::stack<item_t>& stack) {
  item_t item;
  std::stack<item_t> tmp_stack;

  while (!stack.empty()) {
    item = stack.top();
    stack.pop();
    tmp_stack.push(item);
  }

  stack = tmp_stack;
}

void CalcModel::stack_clean(std::stack<item_t>& stack) {
  while (!stack.empty()) {
    stack.pop();
  }
}

int CalcModel::get_precedence(val_t type) {
  int result = 0;

  if (type >= VAL_PLUS && type <= VAL_MINUS) {
    result = 1;
  } else if (type >= VAL_MULTIPLY && type <= VAL_MOD) {
    result = 2;
  } else if (type == VAL_POWER) {
    result = 3;
  } else if (type >= VAL_COS && type <= VAL_LOG) {
    result = 4;
  }

  return result;
}

err_t CalcModel::parse_token(char** src, size_t& index, item_t& item) {
  err_t result = ERR_MALFORMED;
  item.value = 0.0;
  item.type = VAL_UNKNOWN;
  size_t rlen = strlen(*src);

  for (int type = VAL_LBRACKET; (type <= VAL_LOG) && (result == ERR_MALFORMED);
       ++type) {
    if (type == VAL_NUMBER) {
      if (isdigit(**src)) {
        const char* tmp = *src;
        item.value = strtod(*src, src);
        item.type = static_cast<val_t>(type);
        index += *src - tmp;
        result = ERR_SUCCESS;
      }
    } else {
      const char* token_name = op_values[type];
      size_t token_size = strlen(token_name);

      if (rlen >= token_size && strncmp(*src, token_name, token_size) == 0) {
        item.type = static_cast<val_t>(type);
        *src += token_size;
        index += token_size;
        result = ERR_SUCCESS;
      }
    }
  }

  return result;
}

err_t CalcModel::parse_rpn(const char* src) {
  err_t result = ERR_SUCCESS;
  val_t last_token = VAL_UNKNOWN;
  char* ptr = const_cast<char*>(src);
  size_t index = 0;

  while (result == ERR_SUCCESS && *ptr != '\0') {
    if (*ptr == ' ') {
      ptr += 1;
      index += 1;
      continue;
    }

    item_t item = {0.0, VAL_UNKNOWN};

    result = parse_token(&ptr, index, item);

    if (result == ERR_SUCCESS) {
      if ((item.type == VAL_PLUS || item.type == VAL_MINUS) &&
          (last_token == VAL_LBRACKET || last_token == VAL_UNKNOWN)) {
        item_t num = {0.0, VAL_NUMBER};
        result = stack_push(main_stack, num);
        if (result != ERR_SUCCESS) {
          break;
        }
      }

      // Number or variable
      if (item.type == VAL_NUMBER || item.type == VAL_VARIABLE) {
        result = stack_push(main_stack, item);
        // Function or left bracket
      } else if ((item.type >= VAL_COS && item.type <= VAL_LOG) ||
                 item.type == VAL_LBRACKET) {
        result = stack_push(oper_stack, item);
        // Operator
      } else if (item.type >= VAL_PLUS && item.type <= VAL_POWER) {
        while (result == ERR_SUCCESS &&
               stack_peek(oper_stack).type != VAL_UNKNOWN &&
               get_precedence(stack_peek(oper_stack).type) >=
                   get_precedence(item.type)) {
          result = stack_push(main_stack, stack_pop(oper_stack));
        }
        if (result == ERR_SUCCESS) {
          result = stack_push(oper_stack, item);
        }
        // Right bracket
      } else if (item.type == VAL_RBRACKET) {
        while (result == ERR_SUCCESS &&
               stack_peek(oper_stack).type != VAL_UNKNOWN &&
               stack_peek(oper_stack).type != VAL_LBRACKET) {
          result = stack_push(main_stack, stack_pop(oper_stack));
        }
        if (stack_pop(oper_stack).type != VAL_LBRACKET) {
          result = ERR_MISSBRACKET;
        }
        if (result == ERR_SUCCESS && stack_peek(oper_stack).type >= VAL_COS &&
            stack_peek(oper_stack).type <= VAL_LOG) {
          result = stack_push(main_stack, stack_pop(oper_stack));
        }
      }

      last_token = item.type;
    }
  }

  while (result == ERR_SUCCESS && stack_peek(oper_stack).type != VAL_UNKNOWN) {
    item_t item = stack_pop(oper_stack);
    if (item.type == VAL_LBRACKET || item.type == VAL_RBRACKET) {
      result = ERR_MISSBRACKET;
    } else {
      result = stack_push(main_stack, item);
    }
  }

  if (result == ERR_SUCCESS) {
    stack_reverse(main_stack);
  }

  return result;
}

err_t CalcModel::evaluate_rpn() {
  err_t result = ERR_SUCCESS;
  item_t item;

  while (result == ERR_SUCCESS &&
         (item = stack_pop(main_stack)).type != VAL_UNKNOWN) {
    if (item.type == VAL_NUMBER) {
      result = stack_push(oper_stack, item);
    } else if (item.type == VAL_VARIABLE) {
      item.value = x_value;
      result = stack_push(oper_stack, item);
    } else {
      result = math(item);
    }
  }

  if (result == ERR_SUCCESS) {
    if (oper_stack.size() == 1 && stack_peek(oper_stack).type == VAL_NUMBER) {
      out_value = stack_pop(oper_stack).value;
    } else if (oper_stack.size() == 1 &&
               stack_peek(oper_stack).type == VAL_VARIABLE) {
      out_value = x_value;
    } else {
      result = ERR_MALFORMED;
    }
  }

  return result;
}

err_t CalcModel::math(item_t op) {
  err_t result = ERR_SUCCESS;
  item_t first = stack_pop(oper_stack);
  double a, b;

  if (first.type == VAL_UNKNOWN) {
    result = ERR_MALFORMED;
  } else {
    a = first.value;
    if (op.type >= VAL_PLUS && op.type <= VAL_POWER) {
      item_t second = stack_pop(oper_stack);

      if (second.type == VAL_UNKNOWN) {
        result = ERR_MALFORMED;
      } else {
        b = second.value;
      }
    }
  }

  item_t num = {0.0, VAL_NUMBER};

  if (result == ERR_SUCCESS) {
    switch (op.type) {
      case VAL_PLUS:
        num.value = b + a;
        result = stack_push(oper_stack, num);
        break;
      case VAL_MINUS:
        num.value = b - a;
        result = stack_push(oper_stack, num);
        break;
      case VAL_MULTIPLY:
        num.value = b * a;
        result = stack_push(oper_stack, num);
        break;
      case VAL_DIVIDE:
        if (a == 0.0) {
          result = ERR_DIVBYZERO;
        } else {
          num.value = b / a;
          result = stack_push(oper_stack, num);
        }
        break;
      case VAL_MOD:
        if (a == 0.0) {
          result = ERR_DIVBYZERO;
        } else {
          num.value = fmod(b, a);
          result = stack_push(oper_stack, num);
        }
        break;
      case VAL_POWER:
        num.value = pow(b, a);
        result = stack_push(oper_stack, num);
        break;
      case VAL_COS:
        num.value = cos(a);
        result = stack_push(oper_stack, num);
        break;
      case VAL_SIN:
        num.value = sin(a);
        result = stack_push(oper_stack, num);
        break;
      case VAL_TAN:
        num.value = tan(a);
        result = stack_push(oper_stack, num);
        break;
      case VAL_ACOS:
        num.value = acos(a);
        result = stack_push(oper_stack, num);
        break;
      case VAL_ASIN:
        num.value = asin(a);
        result = stack_push(oper_stack, num);
        break;
      case VAL_ATAN:
        num.value = atan(a);
        result = stack_push(oper_stack, num);
        break;
      case VAL_SQRT:
        if (a < 0.0) {
          result = ERR_NEGATIVEROOT;
        } else {
          num.value = sqrt(a);
          result = stack_push(oper_stack, num);
        }
        break;
      case VAL_LN:
        num.value = log(a);
        result = stack_push(oper_stack, num);
        break;
      case VAL_LOG:
        num.value = log10(a);
        result = stack_push(oper_stack, num);
        break;
      default:
        result = ERR_MALFORMED;
        break;
    }
  }

  return result;
}

err_t CalcModel::calculate(const std::string& src, double x) {
  stack_clean(main_stack);
  stack_clean(oper_stack);
  x_value = x;

  err_t result = parse_rpn(src.c_str());

  if (result == ERR_SUCCESS) {
    result = evaluate_rpn();
  }

  return result;
}
