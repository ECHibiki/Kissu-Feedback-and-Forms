package types


// Pretty sure some type inheritence can make this cleaner
type ConfigurationSettings struct{
  DBName string
  DBUserName  string
  DBCredentials  string
  DBAddr  string
  StartupPort string
  SiteName string
  ResoruceDirectory string
}

type ConfigurationInitializerFields struct{
  DBName string
  DBUserName  string
  DBCredentials  string
  DBAddr  string
  ApplicationPassword string
  StartupPort string
  SiteName string
  ResoruceDirectory string
}

type FormDBFields struct{
  ID int64
  FieldJSON string
  UpdatedAt int64
}
type ResponseDBFields struct{
  ID int64
  FK_ID int64
  Identifier string
  ResponseJSON string
  SubmittedAt int64
}
type PasswordsDBFields struct{
  HashedPassword string
  HashSystem string
  HashScrambler string
}
type LoginsDBFields struct{
  TimeAt int64
  Cookie string
  IP string
}
