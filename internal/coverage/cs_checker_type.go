//ff:type feature=coverage type=model lang=csharp
//ff:what C# test coverage checker struct
package coverage

// CsChecker checks C# test coverage using `dotnet test` with the coverlet XPlat
// collector, parsing the resulting Cobertura XML.
type CsChecker struct{}
