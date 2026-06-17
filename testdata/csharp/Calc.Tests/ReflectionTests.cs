namespace Calc.Tests;

using System;
using System.Reflection;
using Xunit;

public class ReflectionTests
{
    [Fact]
    public void ReachesPrivateState()
    {
        var t = typeof(Calculator);

        MethodInfo helper = t.GetMethod("Helper", BindingFlags.NonPublic | BindingFlags.Static);
        helper.Invoke(null, new object[] { 21 });

        FieldInfo total = t.GetField("Total");
        total.SetValue(null, 99);

        PropertyInfo prop = t.GetProperty("Total");
        prop.GetValue(null);
    }
}
