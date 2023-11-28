package uc

import "testing"

const jwtKey = "lidsuaod709j2h*&TUYFTDS"
const jwtTokenTTL = 3600 * 24 * 30 // 30天

func TestJwtDecode(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTAwMDAyLCJuYW1lIjoiIiwiYXZhdGFyIjoiIiwiZXhwIjoxNzAzNjY5MDU5LCJuYmYiOjE3MDEwNzcwNTksImlhdCI6MTcwMTA3NzA1OX0.ZPcSSyerO8IwZLNwJB6FmzwyhSGeDkk8GLuzbot1n3k"
	cli := DefaultJwtEncoder(jwtKey, jwtTokenTTL)
	info, err := cli.decodeJwt(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)
}
