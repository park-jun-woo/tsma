namespace Calc.Tests;

using Xunit;

public class CalculatorTests
{
    [Fact]
    public void Classifies()
    {
        var calc = new Calculator(0);
        Assert.Equal("high", calc.Classify(5, 3));
        Assert.Equal("low", calc.Classify(1, 3));
        Assert.Equal("equal", calc.Classify(3, 3));
    }

    [Fact]
    public void StringArgumentIsNotReflection()
    {
        // The string "GetMethod" here is an argument, not a reflective call,
        // and "MethodInfo" appears only in this comment — neither must fire.
        Assert.False(StringUtils.IsBlank("GetMethod"));
    }
}
