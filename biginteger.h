#define profile(x) std::cerr << x
#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
const int BASE = 10000;
const int LOG = 4;


class BigInteger{
    private:
        mutable std::vector<int> number;
        mutable int sign;
        bool compare(const std::vector<int>& number1, const std::vector<int>& number2) {
            if(number1.size() < number2.size()) {
                return true;
            }
            if(number1.size() > number2.size()) {
                return false;
            }
            size_t counter = 0;
            while(counter < number1.size() && number1[counter] == number2[counter]) {
                ++counter;
            }
            if(counter < number1.size()) {
                return number1[counter] < number2[counter];
            }
            return false;
        }
        void multiply_on_digit(std::vector<int>& num, int digit) {
            
            long long rest = 0;
            int cur_index1 = num.size() - 1;
            long long tmp = 0;
            while(cur_index1 >= 0) {
                tmp = num[cur_index1];
                num[cur_index1] = ((tmp) * (digit) + rest) % BASE;
                rest = ((tmp) * (digit) + rest) / BASE;
                --cur_index1;
            }
            if(rest != 0) {
                num.insert(num.begin(), rest);
            }
        }
        void subtraction(std::vector<int>& number1, std::vector<int>& number2) {
            if(number1 == number2) {
                number1.clear();
                number1[0] = 0;
                return;
            }
            int cur_index1 = number1.size() - 1;
            int cur_index2 = number2.size() - 1;
            bool state = false;
            // std::vector<int>tmp = number1;
            while(cur_index2 >= 0) {
                if(!state) {
                    if(number1[cur_index1] < number2[cur_index2]) {
                        state = true;
                        number1[cur_index1] = BASE + (number1[cur_index1]) - (number2[cur_index2]);
                    } else {
                        number1[cur_index1] = (number1[cur_index1]) - (number2[cur_index2]);
                    }
                } else {
                    if(number1[cur_index1] <= number2[cur_index2]) {
                        number1[cur_index1] = BASE - 1 + (number1[cur_index1]) - (number2[cur_index2]);
                    } else {
                        state = false;
                        number1[cur_index1] = (number1[cur_index1]) - 1 - (number2[cur_index2]);
                    }
                }
                --cur_index2;
                --cur_index1;
            }
            while(cur_index1 >= 0) {
                if(state) {
                    if((number[cur_index1]) == 0) {
                        number1[cur_index1] = BASE - 1;
                    } else {
                        number1[cur_index1] = number1[cur_index1] - 1;
                        state = false;
                    }
                } else {
                    number1[cur_index1] = number1[cur_index1];
                }
                --cur_index1;
            }
            size_t k = 0;
            while(k < number1.size() - 1 && number1[k] == 0) {
                ++k;
            }
            number1.erase(number1.begin(), number1.begin() + k);
        }
        int getDigitAndRest(std::vector<int>& number1, const std::vector<int>& number2) {
            std::vector<int> tmp;
            if(compare(number1, number2)) {
                return 0;
            }
            int l = 1;
            int r = BASE;
            int m = 0;
            while(r > l) {
                m = (l + r) / 2;
                tmp = number2;
                multiply_on_digit(tmp, m);
                if(compare(number1, tmp)) {
                    r = m;
                } else {
                    l = m + 1;
                }
            }
            tmp = number2;
            multiply_on_digit(tmp, r - 1);
            subtraction(number1, tmp);
            return r - 1;
        }
    public:
        BigInteger(int n) {
            if(n == 0) {
                number.push_back(0);
                sign = 1;
            } else {
                if(n < 0) {
                    sign = -1;
                } else {
                    sign = 1;
                }
                n = std::abs(n);
                while(n > 0) {
                    number.push_back(n % BASE);
                    n /= BASE;
                }
                std::reverse(number.begin(), number.end());
            }
        }
        BigInteger() {}
        BigInteger& operator=(const BigInteger& bigint) {
            number = bigint.number;
            sign = bigint.sign;
            return *this;
        }
        BigInteger& operator= (BigInteger& bigint) {
            number = bigint.number;
            sign = bigint.sign;
            return *this;
        }
        BigInteger& operator=(int n) {
            number.clear();

            if(n == 0) {
                number.push_back(0);
                sign = 1;
                return *this;
            }
            if(n < 0) {
                sign = -1;
            } else {
                sign = 1;
            }
            n = std::abs(n);
            while(n > 0) {
                number.push_back(n % BASE);
                n /= BASE;
            }
            std::reverse(number.begin(), number.end());
            return *this;
        }
        BigInteger& operator+= (const BigInteger& bigint) {
            if(sign == -1 && bigint.sign == -1) {
                // *this *= (-1);
                sign *= -1;
                bigint.sign *= -1;
                *this += bigint;
                sign *= -1;
                bigint.sign *= -1;
                return *this;
            }
            if(sign != -1 && bigint.sign == -1) {
                bigint.sign *= -1;
                *this -= bigint;
                bigint.sign *= -1;
                return *this;
            }
            if(sign == -1 && bigint.sign != -1) {
                // *this *= (-1);
                sign *= -1;
                *this -= bigint;
                // *this *= (-1);
                sign *= -1;
                return *this;
            }
            if(bigint == 0) {
                return *this;
            }
            long long rest = 0;
            int cur_index1 = static_cast<int> (number.size() - 1);
            int cur_index2 = static_cast<int> (bigint.number.size() - 1);
            long long tmp = 0;
            if(cur_index1 >= cur_index2) {
                while(cur_index2 >= 0) {
                    tmp = number[cur_index1];
                    number[cur_index1] = ((tmp) + (bigint.number[cur_index2]) + rest) % BASE;
                    rest = ((tmp) + (bigint.number[cur_index2]) + rest) / BASE;
                    --cur_index1;
                    --cur_index2;
                }
                while(cur_index1 >= 0) {
                    tmp = number[cur_index1];   
                    number[cur_index1] = (tmp + rest) % BASE;
                    rest = (tmp + rest) / BASE;
                    --cur_index1;
                }
                if(rest != 0) {
                    number.insert(number.begin(), rest);
                }
            } else {
                std::reverse(number.begin(), number.end());
                cur_index1 = 0;
                int sizee = static_cast<int>(number.size());
                while(cur_index1 < sizee) {
                    tmp = number[cur_index1];
                    number[cur_index1] = (tmp + bigint.number[cur_index2] + rest) % BASE;
                    rest = (tmp + bigint.number[cur_index2] + rest) / BASE;
                    --cur_index2;
                    ++cur_index1;
                }
                while(cur_index2 >= 0) {
                    number.push_back((bigint.number[cur_index2] + rest) % BASE);
                    rest = (rest + bigint.number[cur_index2]) / BASE;
                    --cur_index2;
                }
                if(rest != 0) {
                    number.push_back(rest);
                }
                std::reverse(number.begin(), number.end());
                
            }
            return *this;
        }
        BigInteger& operator-=(const BigInteger&);
        
        BigInteger& operator*= (const BigInteger& bigint) {
            std::vector<int> ans (bigint.number.size() + number.size());
	// profile("*=");
            std::reverse(number.begin(), number.end());
            std::reverse(bigint.number.begin(), bigint.number.end());
            long long rest = 0;
            long long cur_number = 0;
            for(size_t i = 0; i < number.size(); ++i) {
                for(size_t j = 0; j < bigint.number.size() || rest != 0; ++j) {
                    cur_number = ans[i + j] + number[i] * (j < bigint.number.size() ? bigint.number[j] : 0) + rest;
                    ans[i + j] = cur_number % BASE;
                    rest = cur_number / BASE;
                }
            }
            sign *= bigint.sign;
            size_t k = ans.size() - 1;
            while(ans.size() > 1 && ans[k] == 0) {
                ans.pop_back();
                --k;
            }
            std::reverse(bigint.number.begin(), bigint.number.end());
            std::reverse(ans.begin(), ans.end());
            number = ans;
            return *this;
        }
        BigInteger& operator/= (const BigInteger& bigint) {
            // profile("/=");
            if(bigint == 1) {
	
                return *this;
            }
            if(number.size() == 1 && number[0] == 0) {
                return *this;
            }
            sign *= bigint.sign;
            std::vector<int>cur_number;
            size_t index = 0;
            while(index < number.size() && cur_number.size() < bigint.number.size()) {
                cur_number.push_back(number[index]);
                ++index;
            }
            if(compare(cur_number, bigint.number)) {
                if(index == number.size()) {
                    *this = 0;
                    return *this;
                }
                cur_number.push_back(number[index]);
                ++index;
            }
            int digit = 0;
            std::vector<int> ans;
            int temp = 0;
            while(index < number.size()) {
                temp = number[index];
                digit = getDigitAndRest(cur_number, bigint.number);  
                ans.push_back(digit);
                // number[index] = digit;
                if(cur_number.size() == 1 && cur_number[0] == 0) {
                    cur_number.clear();
                }
                cur_number.push_back(temp);
                ++index;
            }
            digit = getDigitAndRest(cur_number, bigint.number);
            ans.push_back(digit);
            number = ans;
            return *this;

        }
        BigInteger& operator%= (const BigInteger&);
        BigInteger& operator++() {
            *this += 1;

            return *this;
        }
        BigInteger& operator--() {
            *this -= 1;
	
            return *this;
        }
        BigInteger operator-() const {
            if(*this == 0) {

                BigInteger bigint;
                bigint.number = number;
                bigint.sign = 1;
                return bigint;
            }
            BigInteger bigint;
            if(sign < 0) {
                bigint.number.clear();
                bigint.sign = 1;
                for(size_t i = 0; i < number.size(); ++i) {
                    bigint.number.push_back(number[i]);
                }
                return bigint;
            } 
            bigint.sign = -1;
            bigint.number.clear();
            for(size_t i = 0; i < number.size(); ++i) {
                bigint.number.push_back(number[i]);
            }
            return bigint;
        }
        BigInteger& operator= (std::string& s) {
            number.clear();
            if(s.size() == 0) {
                number.push_back(0);
                sign = 1;
                return *this;
            }
            if(s.size() > 100) {
	// profile(s);
            }
            std::string tmp = "";
            if(s[0] == '-') {
                s.erase(s.begin());
                sign = -1;
            } else {
                sign = 1;
            }
            size_t rest = s.size() % LOG;
            if(rest != 0) {
                for(size_t i = 0; i < rest; ++i) {
                    tmp += s[i];
                }
                number.push_back(std::stoi(tmp));
            }
            tmp = "";
            if(rest < s.size()) {
                for(size_t i = rest; i + LOG <= s.size(); i += LOG) {
                    tmp = "";
                    for(size_t j = i; j < i + LOG; ++j) {
                        tmp += s[j];
                    }
                    number.push_back(std::stoi(tmp));
                }
            }
            return *this;
        }
        BigInteger operator++(int) {
            BigInteger ans = *this;
	
            *this += 1;
            return ans;
        }
        BigInteger operator--(int) {
            BigInteger ans = *this;
	
            *this -= 1;
            return ans;
        }
        explicit operator bool() {
            return !(number.size() == 1 && number[0] == 0);
        }
        std::string toDecimal() const {
            std::string ans = "";
            if(sign == -1) {
                ans += "-";
            }
            ans += std::to_string(number[0]);
            std::string temp = "";
            for(size_t i = 1; i < number.size(); ++i) {
                temp = std::to_string(number[i]);
                std::reverse(temp.begin(), temp.end());
                while(temp.size() < LOG) {
                    temp.push_back('0');
                }
                std::reverse(temp.begin(), temp.end());
                ans += temp;
            }
            return ans;

        }
        std::string toString() const {
            return this->toDecimal();
        }
        friend std::ostream& operator<< (std::ostream&, const BigInteger&);
        friend std::istream& operator>> (std::istream&, BigInteger&);
        friend bool operator<(const BigInteger&, const BigInteger&);
        friend bool operator==(const BigInteger&, const BigInteger&);
};


bool operator<(const BigInteger& bigint1, const BigInteger& bigint2) {
    if(bigint1.sign == -1 && bigint2.sign != -1) {
	
        return true;
    }
    if(bigint1.sign == -1 && bigint2.sign == -1) {
        return ((-bigint2) < (-(bigint1)));
    }
    if(bigint1.sign != -1 && bigint2.sign == -1) {
        return false;
    }
    if(bigint1.number.size() < bigint2.number.size()) {
        return true;
    }
    if(bigint1.number.size() > bigint2.number.size()) {
        return false;
    }
    size_t counter = 0;
    while(counter < bigint1.number.size() && bigint1.number[counter] == bigint2.number[counter]) {
        ++counter;
    }
    if(counter < bigint1.number.size()) {
        return bigint1.number[counter] < bigint2.number[counter];
    }
    return false;
}


bool operator>(const BigInteger& bigint1, const BigInteger& bigint2) {
    return (bigint2 < bigint1);
	
}


bool operator==(const BigInteger& bigint1, const BigInteger& bigint2){
    size_t i = 0;
	// profile("==");
    if(bigint1.sign != bigint2.sign) {
        return false;
    }
    if(bigint1.number.size() != bigint2.number.size()) {
        return false;
    }
    while(i < bigint1.number.size() && bigint1.number[i] == bigint2.number[i]) {
        ++i;
    }
    if(i == bigint1.number.size())  {
        return true;
    }
    return false;
}


bool operator<= (const BigInteger& bigint1, const BigInteger& bigint2){
    return (bigint2 == bigint1 || bigint1 < bigint2);
	
}


bool operator>= (const BigInteger& bigint1, const BigInteger& bigint2) {
    return (bigint1 > bigint2 || bigint1 == bigint2);
	
}


bool operator!= (const BigInteger& bigint1, const BigInteger& bigint2) {
    return (!(bigint1 == bigint2));
	
}


BigInteger operator+ (const BigInteger& bigint1, const BigInteger& bigint2) {
    BigInteger ans = bigint1;
	
    ans += bigint2;
    return ans;
}


BigInteger operator* (const BigInteger& bigint1, const BigInteger& bigint2) {
    BigInteger ans = bigint1;
	
    ans *= bigint2;
    return ans;
}


BigInteger operator- (const BigInteger& bigint1, const BigInteger& bigint2) {
    BigInteger ans = bigint1;
	
    ans -= bigint2;
    return ans;
}


BigInteger operator/ (const BigInteger& bigint1, const BigInteger& bigint2) {
    BigInteger ans = bigint1;
	
    ans /= bigint2;
    return ans;
}



BigInteger operator% (const BigInteger& bigint1, const BigInteger& bigint2) {
    BigInteger ans = bigint1;
	
    ans %= bigint2;
    return ans;
}


void substr_vectors(std::vector<int>& n1, const std::vector<int>& n2) {
    std::reverse(n1.begin(), n1.end());
    int cur_index1 = 0;
    int sizee = static_cast<int> (n1.size());
    int cur_index2 = n2.size() - 1;
    bool state = false;
    while(cur_index1 < sizee) {
        if(!state) {
            if(n1[cur_index1] > n2[cur_index2]) {
                state = true;
                n1[cur_index1] = BASE + (n2[cur_index2]) - (n1[cur_index1]);
            } else {
                n1[cur_index1] = (n2[cur_index2]) - (n1[cur_index1]);
            }
        } else {
            if(n2[cur_index2] <= n1[cur_index1]) {
                n1[cur_index1] = BASE - 1 + (n2[cur_index2]) - (n1[cur_index1]);
            } else {
                state = false;
                n1[cur_index1] = (n2[cur_index2]) - 1 - (n1[cur_index1]);
                
            }
        }
        --cur_index2;
        ++cur_index1;
    }
    while(cur_index2 >= 0) {
        if(state) {
            if((n2[cur_index2]) == 0) {
                n1.push_back(BASE - 1);
            } else {
                n1.push_back(n2[cur_index2] - 1);
                state = false;
            }
        } else {
            n1.push_back(n2[cur_index2]);
        }
        --cur_index2;
    }
    size_t k = n1.size() - 1;
    while(n1.size() > 1 && n1[k] == 0) {
        --k;
        n1.pop_back();
    }
    std::reverse(n1.begin(), n1.end());
}


BigInteger& BigInteger::operator-= (const BigInteger& bigint) {
    if(*this == bigint) {
        number.clear();
        number.push_back(0);
        sign = 1;
        return *this;
    }
    if(*this < 0 && (0 < bigint || bigint == 0)) {
        sign *= (-1);
        *this += bigint;
        sign *= (-1);
        return *this;
    }
    if(*this < 0 && bigint < 0) {
        sign *= (-1);
        bigint.sign *= (-1);
        *this -= bigint;
        sign *= (-1);
        bigint.sign *= (-1);
        return *this;
    }
    if((0 < *this || *this == 0) && bigint < 0) {
        ////////////
        bigint.sign *= (-1);
        *this += (bigint);
        bigint.sign *= (-1);
        return *this;
    }
    if(*this < bigint) {
        // profile("test: " + this->toDecimal() + " " + bigint.toDecimal())
        substr_vectors(number, bigint.number);
        sign *= (-1);
        return *this;
    }
    int cur_index1 = number.size() - 1;
    int cur_index2 = bigint.number.size() - 1;
    bool state = false;
    while(cur_index2 >= 0) {
        if(!state) {
            if(number[cur_index1] < bigint.number[cur_index2]) {
                state = true;
                number[cur_index1] = BASE + (number[cur_index1]) - (bigint.number[cur_index2]);
            } else {
                number[cur_index1] = (number[cur_index1]) - (bigint.number[cur_index2]);
            }
        } else {
            if(number[cur_index1] <= bigint.number[cur_index2]) {
                number[cur_index1] = BASE - 1 + (number[cur_index1]) - (bigint.number[cur_index2]);
            } else {
                state = false;
                number[cur_index1] = (number[cur_index1]) - 1 - (bigint.number[cur_index2]);
            }
        }
        --cur_index2;
        --cur_index1;
    }
    while(cur_index1 >= 0) {
        if(state) {
            if((number[cur_index1]) == 0) {
                number[cur_index1] = BASE - 1;
            } else {
                number[cur_index1] = number[cur_index1] - 1;
                state = false;
            }
        } else {
            number[cur_index1] = number[cur_index1];
        }
        --cur_index1;
    }
    size_t k = 0;
    while(k < number.size() - 1 && number[k] == 0) {
        ++k;
    }
    number.erase(number.begin(), number.begin() + k);
    return *this;
}


BigInteger& BigInteger::operator%= (const BigInteger& bigint) {
    BigInteger bg = *this;
    bg /= bigint;
    bg *= bigint;
    *this -= bg;
    return *this;
}


std::istream& operator>>(std::istream& in, BigInteger& bigint) {
    char ch = '\0';
    std::string ans = "";

    while(std::isspace(ch = in.get()));

    in.unget();

    while(true) {
        in.get(ch);
        if(in.eof() || std::isspace(ch)) {
            break;
        }
        ans += ch;
    }
    bigint = ans;
    return in;
}


std::ostream& operator<< (std::ostream& out, const BigInteger& bg) {
    out << bg.toString();
    return out;
}


class Rational {
    private:
        BigInteger numerator;
        BigInteger denominator;
    public:
        Rational() {
            numerator = 0;
            denominator = 1;
        }
        Rational (int n) {
            numerator = n;
            denominator = 1;
        }
        Rational(const BigInteger& bigint) {
            numerator = bigint;
            denominator = 1;
        }
        Rational (const Rational& r) {
            numerator = r.numerator;
            denominator = r.denominator;
        }
        Rational& operator= (int n) {
            numerator = n;
            denominator = 1;
            return *this;
        }
        Rational& operator= (const Rational& fraction) {
            numerator = fraction.numerator;
            denominator = fraction.denominator;
            return *this;
        }
        BigInteger EuclidAlgorithm(BigInteger bigint1, BigInteger bigint2) {
            if(bigint1 < 0 && bigint2 > 0) {
                bigint1 *= (-1);
                return EuclidAlgorithm(bigint1, bigint2);
            } 
            if(bigint1 > 0 && bigint2 < 0) {
                bigint2 *= (-1);
                return -EuclidAlgorithm(bigint1, bigint2);
            }
            if(bigint1 < 0 && bigint2 < 0) {
                bigint1 *= (-1);
                bigint2 *= (-1);
                return -EuclidAlgorithm(bigint1, bigint2);
            }
            while(bigint2 != 0 && bigint1 != 0) {
                if(bigint1 >= bigint2) {
                    bigint1 %= bigint2;
                } else {
                    bigint2 %= bigint1;
                }
            }
            if(bigint1 == 0) {
                return bigint2;
            }
            return bigint1;
        }
        Rational& operator+= (const Rational& fraction) {
            numerator *= fraction.denominator;
            numerator += (denominator * fraction.numerator);
            denominator *= fraction.denominator;
            BigInteger gcd = EuclidAlgorithm(numerator, denominator);
            numerator /= gcd;
            denominator /= gcd;
            return *this;
        }
        Rational& operator-= (const Rational& fraction) {
            numerator *= fraction.denominator;
            numerator -= (fraction.numerator * denominator);
            denominator *= fraction.denominator;
            BigInteger gcd = EuclidAlgorithm(numerator, denominator);
            numerator /= gcd;
            denominator /= gcd;
            return *this;
        }
        Rational& operator*= (const Rational& fraction) {
            numerator *= fraction.numerator;
            denominator *= fraction.denominator;
            BigInteger gcd = EuclidAlgorithm(numerator, denominator);
            numerator /= gcd;
            denominator /= gcd;
            return *this;
        }
        Rational& operator/= (const Rational& fraction) {
            numerator *= fraction.denominator;
            denominator *= fraction.numerator;
            BigInteger gcd = EuclidAlgorithm(numerator, denominator);
            numerator /= gcd;
            denominator /= gcd;
            return *this;
        }
        Rational operator-() const {
            Rational ans = *this;
            ans.numerator *= (-1);
            return ans;
        }
        std::string toString() const {
            if(denominator == 1) {
                return numerator.toDecimal();
            }
            return (numerator.toDecimal() + "/" + denominator.toDecimal());
        }
        std::string asDecimal(size_t precision = 0) {
            BigInteger bg = numerator;
            std::string ans = "";
            if(numerator < 0) {
                ans += "-";
                bg *= -1;
            }
            ans += ((bg / denominator).toDecimal());
            if(precision == 0) {
                return ans;
            }
            ans += ".";
            bg %= denominator;
            for(size_t i = 0; i < precision; ++i) {
                bg *= 10;
                if(bg < denominator) {
                    ans += '0';
                }
            }
            ans += (bg / denominator).toString();
            return ans;
        }
        explicit operator double() {
            double ans = std::atof(asDecimal(20).c_str());
            return ans;
        }
        friend bool operator< (const Rational&, const Rational&);
        friend bool operator==(const Rational&, const Rational&);
        friend std::ostream& operator<<(std::ostream&, const Rational&);
};


std::ostream& operator<<(std::ostream& out, const Rational& fraction) {
    out << fraction.toString();
    return out;
}


bool operator< (const Rational& fraction1, const Rational& fraction2) {
    return (fraction1.numerator * fraction2.denominator < fraction2.numerator * fraction2.denominator);
}


bool operator>(const Rational& fraction1, const Rational& fraction2) {
    return (fraction2 < fraction1);
}


bool operator==(const Rational& fraction1, const Rational& fraction2) {
    return (fraction1.numerator == fraction2.numerator && fraction1.denominator == fraction2.denominator);
}


bool operator<= (const Rational& fraction1, const Rational& fraction2) {
    return (fraction1 < fraction2 || fraction1 == fraction2);
}


bool operator>=(const Rational& fraction1, const Rational& fraction2) {
    return (fraction2 < fraction1 || fraction2 == fraction1);
}


bool operator!=(const Rational& fraction1, const Rational& fraction2) {
    return !(fraction1 == fraction2);
}


Rational operator+ (const Rational& frac1, const Rational& frac2) {
    Rational ans = frac1;
    ans += frac2;
    return ans;
}


Rational operator- (const Rational& frac1, const Rational& frac2) {
    Rational ans = frac1;
    ans -= frac2;
    return ans;
}


Rational operator* (const Rational& frac1, const Rational& frac2) {
    Rational ans = frac1;
    ans *= frac2;
    return ans;
}


Rational operator/ (const Rational& frac1, const Rational& frac2) {
    Rational ans = frac1;
    ans /= frac2;
    return ans;
}