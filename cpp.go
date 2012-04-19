package kodetex

import (
  "log"
  "os"
)

const (
  emptyFlag       byte = 0
  symbolFlag      byte = 1
  starFlag        byte = 1 << 1
  nonGreedyFlag   byte = 1 << 2
  beginGroupFlag  byte = 1 << 3
  endGroupFlag    byte = 1 << 4

  RuleEndsWithEscape = "The rule ends with an escape symbol"
)

type CppParser struct {
  File *os.File;
};

type Exp struct {
  backtrack []uint;
  flags []byte;
  symbols []byte;
};

type Matcher interface {
  Matches(string, uint);
};

type Group struct {
  children []Matcher;
};

func (*CppParser) Parse() {

  grammar := "a(c/d)/e";

//  grammar := "\\/\\*.*?\\*\\//\\/\\/.*?\n";
  length := uint(len(grammar));

  e := &Exp {
    backtrack: make([]uint, 0, length),
    flags: make([]byte, 0, length),
    symbols: make([]byte, 0, length),
  }

  var i uint;
  var symbol byte;

  for i = 0 ; i < length ; i++ {

    symbol = grammar[i];

    if symbol != '\\' {
      e.parse(symbol);
      continue;
    }

    i++;
    if i >= length {
      panic(RuleEndsWithEscape);
    }
    e.append(grammar[i]);

  }

  log.Println(string(e.symbols));
  log.Println(e.flags);
  log.Println(e.backtrack);
//  log.Println(e.matches("abbbbbbcd"));
}

func (e *Exp) parse(symbol byte) {

  switch symbol {
    case '*':
      e.apply(starFlag);
    case '?':
      e.apply(nonGreedyFlag);
    case '(':
      e.apply(beginGroupFlag);
    case ')':
      e.advance();
      e.apply(endGroupFlag);
    case '/':
      e.updateNext();
    case '.':
      e.advance();
    default:
      e.append(symbol);
  }

}

func (e *Exp) append(symbol byte) {

  e.symbols = append(e.symbols, symbol);
  e.flags = append(e.flags, symbolFlag);
  e.backtrack = append(e.backtrack, 0);

}

func (e *Exp) advance() {

  e.flags = append(e.flags, emptyFlag);
  e.backtrack = append(e.backtrack, 0);

}

func (e *Exp) apply(flag byte) {

  e.flags[len(e.flags)-1] |= flag;

}

func (e *Exp) updateNext() {

  length := uint(len(e.backtrack));
  if length == 0 {
    return;
  }

  var i uint = length - 1;
  lastValue := e.backtrack[i];
  for ; e.backtrack[i] == lastValue && !e.isBeginGroup(i) ; i-- {
    e.backtrack[i] = length;
    if i == 0 {
      break;
    }
  }

}

func (e *Exp) isSymbol(index uint) bool {
  return e.flags[index] & symbolFlag == symbolFlag;
}

func (e *Exp) isStar(index uint) bool {
  return e.flags[index] & starFlag == starFlag;
}

func (e *Exp) isBeginGroup(index uint) bool {
  return e.flags[index] & beginGroupFlag == beginGroupFlag;
}

func (e *Exp) length() uint {
  return uint(len(e.flags));
}

func (e *Exp) matches(text string) uint {
  var j uint = 0;
  var i uint = 0;
  var textLength uint = uint(len(text));
  var expLength uint = e.length();

  var isStar bool;

  for ; i < textLength && j < expLength ; {

    isStar = e.isStar(j);

    if e.isSymbol(j) {
      if e.symbols[j] != text[i] {
        if isStar {
          j++;
          continue;
        }
        if e.backtrack[j] == 0 {
          return 1;
        }
        j = e.backtrack[j];
        continue;
      }
      i++;
      if !isStar {
        j++;
      }
      continue;
    }
  }
  return 0;
}
