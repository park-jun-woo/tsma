namespace Calc.Tests;

using Xunit;

public class StringUtilsTests
{
    [Fact]
    public void Repeats()
    {
        Assert.Equal("aa", StringUtils.Repeat("a", 2));
    }
}
