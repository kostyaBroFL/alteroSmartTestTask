package flagenv

//
// import (
//     "flag"
//     "github.com/stretchr/testify/suite"
//     "os"
//     "testing"
// )
//
// type envFlagParserTest struct {
//     suite.Suite
// }
//
// func (p *envFlagParserTest) TestMustParseInt(t *testing.T) {
//     envName := "TEST_INT"
//     os.Setenv(envName, "4")
//     flagV := flag.Int("test_int", 0, "description")
//     p.Assert().Equal(4, mustParseInt(flagV, envName))
// }
//
// func TestEnvFlagParser(t *testing.T) {
//     suite.Run(t, new(envFlagParserTest))
// }
