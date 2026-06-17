namespace Calc;

public static class StringUtils
{
    public static bool IsBlank(string s)
    {
        if (s == null)
        {
            return true;
        }
        return s.Trim().Length == 0;
    }

    public static string Repeat(string s, int times)
    {
        var result = "";
        for (var i = 0; i < times; i++)
        {
            result += s;
        }
        return result;
    }
}
