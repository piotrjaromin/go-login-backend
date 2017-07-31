package accounts

import (
        . "github.com/smartystreets/goconvey/convey"
        "testing"
)

func TestEncrypt(t *testing.T) {

        encrypt := CreateEncrypt()

        Convey("For hash ", t, func() {

                Convey("should return hashed password and salt", func() {

                        hashed, salt := encrypt.Hash("pass")

                        So(string(hashed), ShouldNotBeBlank)
                        So(salt, ShouldNotBeBlank)
                })

        })

        Convey("For validate", t, func() {

                salt := "testSalt"
                pass := Password("testPass")
                okHashed := Password("d71b6e68a1bddecc0451569b830e7e6d2821d1b3512e611684cec6870ed49d23895ff8648d7b4bcb48deeb807e6ef4430fb53a2f8090b04ada90fd0dbc472069")

                Convey("Should return ok for valid password", func() {

                        ok := encrypt.Validate(pass, okHashed, salt)

                        So(ok, ShouldBeTrue)
                })

                Convey("Should detect invalid password", func() {
                        ok := encrypt.Validate(pass, Password("notOk"), salt)

                        So(ok, ShouldBeFalse)
                })
        })
}
