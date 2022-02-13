#include"type.h"
#include<iostream>
#include<Windows.h>
using namespace std;
typedef void(*func)(void);
//using func = void(*)();
int main()
{
    cout << "hello debug" << endl;
    HMODULE hDll = LoadLibrary(L"analyser.dll");
    auto fn = (func) GetProcAddress(hDll, "Debug");
    if (fn != nullptr)
    {
        cout << "start debug" << endl;
        fn();
    }
    FreeLibrary(hDll);
    cout << "end debug" << endl;
}