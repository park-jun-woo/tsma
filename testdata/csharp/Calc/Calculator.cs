namespace Calc;

using System;

public partial class Calculator
{
    public int Total { get; set; }

    [Obsolete("demo")]
    public Calculator(int seed)
    {
        Total = seed;
    }

    public string Classify(
        int score,
        int threshold)
    {
        if (score > threshold)
        {
            return "high";
        }
        if (score < threshold)
        {
            return "low";
        }
        return "equal";
    }

    private static int Helper(int n)
    {
        return n * 2;
    }

    public class Inner
    {
        public void Ping()
        {
        }
    }
}
