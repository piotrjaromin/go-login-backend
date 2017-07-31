package jwtTokens

import (
	"strings"
        . "github.com/smartystreets/goconvey/convey"
		"testing"
)

func TestService(t *testing.T) {

	const siginKey = "someTestKey"
	servce := Create(siginKey)
	user := "testUser"
	userID := "testUserId"

	Convey("Service should", t, func() {

		Convey("create valid token", func() {
			token, err := servce.GenerateToken(user, userID)
			So(token, ShouldNotBeBlank)
			So(err, ShouldBeNil)

			tokenElements := strings.Split(token, ".")
			So(tokenElements, ShouldHaveLength, 3)

			So(servce.Validate(token, "username", user), ShouldBeTrue)
			So(servce.Validate(token, "userId", userID), ShouldBeTrue)
		})

		Convey("return false for invalid token", func() {
			So(servce.Validate("randomToken", "username", user), ShouldBeFalse)
			So(servce.Validate("randomToken", "userId", userID), ShouldBeFalse)
		})
	})

}