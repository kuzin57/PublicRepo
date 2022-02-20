#include <iostream>
#include <cstring>
#include <string.h>
#pragma once


class String {
    private:
        char* str; 
        size_t capacity;
        size_t size; 
    public:
        String(const char* s) {
            size = strlen(s);
            size_t len = size * 2;
            str = new char[len];
            capacity = len;
            memcpy(str, s, size); 
            str[size] = '\0';
        }
        String(size_t n, char ch) {
            size_t len = 2 * n;
            size = n;
            str = new char[len];
            memset(str, ch, n);
            str[n] = '\0';
        }
        String() {
            str = new char[16];
            str[0] = '\0';
            capacity = 16;
            size = 0;
        }
        String(const String& s){
            size_t len = s.length() * 2;
            size = s.length();
            str = new char[len];
            capacity = len;
            memcpy(str, s.str, size);
            str[size] = '\0';
        }
        explicit String(char ch) {
            capacity = 16;
            str = new char[capacity];
            str[0] = ch;
            str[1] = '\0';
            size = 1;
        }
        ~String() {
           delete[] str;
        }
        String& operator=(const String& s) {
            if (s.size > capacity) {
                delete[] str;
                str = new char[s.size * 2];
            }
            size = s.length();
            memcpy(str, s.str, size);
            str[size] = '\0';
            return *this;
        }
        size_t length() const {
            return size;
        }
        void push_back(char ch) {
            if (size + 1 > capacity) {
                reallocate_second();
            }
            str[size] = ch;
            ++size;
            str[size] = '\0';
        }
        char pop_back() { 
            char last = str[size - 1]; 
            str[size - 1] = '\0'; 
            --size; 
            if (size * 4 < capacity) {
                char* tmp = new char[size + 1];
                memcpy(tmp, str, size);
                delete[] str;
                str = new char[capacity / 2];
                capacity /= 2;
                memcpy(str, tmp, size);
                str[size] = '\0';
                delete[] tmp; 
            }
            return last;
        }
        char front() const {
            return str[0]; 
        }
        char& front() {
            return str[0];
        }
        char& back() {
            return str[size - 1];
        }
        char back() const { 
            return str[size - 1];
        }
        void reallocate_second() {
            size_t len = size * 2 + 1;
            char* str1 = new char[size + 1];
            str1[0] = '\0';
            memcpy(str1, str, size);
            delete[] str;
            str = new char[len];
            capacity = len;
            memcpy(str, str1, size);
            str[size] = '\0';
            delete[] str1;
        }
        void reallocate(const String& s) {
            capacity = (size + s.length()) * 2;
            char* str1 = new char[size + 1]; 
            // str1[0] = '\0';
                
            memcpy(str1, str, size);
            str1[size] = '\0';
            delete[] str;
            str = new char[capacity];
            memcpy(str, str1, size);
            str[size] = '\0';
            delete[] str1;
        }
        String& operator+=(const String& s) { 
            size_t str_len = size;
            if (capacity <= str_len + s.length()) {
                reallocate(s);
            }
            char* tmp = new char[s.length() + 1];
            tmp[0] = '\0';
            memcpy(tmp, s.str, s.size);
            tmp[s.length()] = '\0';
            strcat(str, tmp);
            size += s.size;
            delete[] tmp;
            return *this;
        }
        char operator[] (size_t index) const {
            return str[index];
        }
        char& operator[] (size_t index) {
            return str[index];
        }
        String& operator+= (char ch) {
            push_back(ch);
            return *this;
        }
        String operator+(String& s) {
            String ans;
            ans += *this;
            ans += s;
            return ans;
        }
        String operator+(char ch) {
            String ans;
            ans += *this;
            ans += ch;
            return ans;
        }
        size_t find(const String& substring) const{
            size_t len = substring.length();
            bool state = false;
            size_t cur_index = 0;
            size_t begin_index = 0;
            for (size_t i = 0; i < size; ++i) {
                if (str[i] == substring[cur_index] && !state) {
                    begin_index = i;
                    if(cur_index + 1 == len) {
                        return begin_index;
                    }
                    state = true;
                    cur_index++;
                    continue;
                }
                if (str[i] != substring[cur_index] && state) {
                    begin_index = 0;
                    cur_index = 0;
                    state = false;
                    continue;
                }
                if (str[i] == substring[cur_index] && state) {
                    if(cur_index + 1 == len) {
                        return begin_index;
                    }
                    cur_index++;
                    continue;
                }
            }
            return size;
        }
        size_t rfind(const String& substring) const {
            bool state = false;
            size_t len1 = substring.length();
            size_t cur_index = len1 - 1;
            size_t cur_length = 0;
            for (size_t i = size - 1; ; --i) {
                if (str[i] == substring[cur_index] && !state) {
                    if (cur_length + 1 == len1) {
                        return i;
                    } 
                    state = true;
                    cur_index--;
                    cur_length = 1;
                    continue;
                }
                if (str[i] == substring[cur_index] && state) {
                    if(cur_length + 1 == len1) {
                        return i;
                    }
                    cur_index--;
                    cur_length++;
                }
                if (str[i] != substring[cur_index] && state) {
                    state = false;
                    cur_length = 0;
                    cur_index = len1 - 1;
                    continue;
                }
            }
            return size;
        }
        String substr(size_t start, size_t count) const{
            char* str1 = new char[count + 1];
            for (size_t i = start; i < start + count; ++i) {
                str1[i - start] = str[i];
            }
            str1[count] = '\0';
            String ss = str1;
            delete[] str1;
            return ss;
        }
        void clear() {
            delete[] str;
            str = new char[16];
            str[0] = '\0';
            capacity = 16;
            size = 0;
        }
        bool operator==(const String& s) const {
            if (s.size != size) {
                return false;
            }
            return strcmp(str, s.str) == 0;
        }
        bool empty() { 
            return (size == 0); 
        }
        friend std::ostream& operator<< (std::ostream&, const String&);
        friend std::istream& operator>> (std::istream&, String&);
        // вот это для выражений типа char + String 
        friend String operator+ (char, String&);
};


std::ostream& operator << (std::ostream& out, const String& s) {
    for (size_t i = 0; i < s.size; ++i) {
        out << s[i];
    }
    return out;
}


std::istream& operator >> (std::istream& in, String& s) {
    s.clear();
    char ch = 'a';
    while (!isspace(ch)) {
        in.get(ch);
        if (in.eof()) {
            break;
        }
        if (!isspace(ch)) {
            s += ch;
        }
    }
    return in;
}


String operator+ (char ch, String& s) {
    String s1(1, ch);
    s1 += s;
    return s1;
}
