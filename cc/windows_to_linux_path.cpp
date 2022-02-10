#include <iostream>
#define MAX 1024
using namespace std;
void convert(char* str)
{
    int i = 3;
    cout << "\n" << "/" << str[0] << "/";
    while (i < MAX)
    {
        if (str[i] == '\0') { cout << endl; return; };
        if (str[i] == '\\') cout << '/'; else cout << str[i];
        i++;
    }
}


int   main(int   argc, char* argv[])
{
    if (argc < 2)
    {
        return -1;
    }
    try
    {
        if (argv[1][1] == ':')
        {
            convert(argv[1]);
        }
    } catch (...)
    {
        cout << "参数错误" << endl;
    }
    return   0;
}


