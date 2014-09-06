package miniredis

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestAuth(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c, err := redis.Dial("tcp", s.Addr())
	ok(t, err)

	// We accept all AUTH
	_, err = c.Do("AUTH", "foo", "bar")
	ok(t, err)
}

func TestEcho(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c, err := redis.Dial("tcp", s.Addr())
	ok(t, err)

	r, err := redis.String(c.Do("ECHO", "hello\nworld"))
	ok(t, err)
	equals(t, "hello\nworld", r)
}

func TestSelect(t *testing.T) {
	s, err := Run()
	ok(t, err)
	defer s.Close()
	c, err := redis.Dial("tcp", s.Addr())
	ok(t, err)

	_, err = redis.String(c.Do("SET", "foo", "bar"))
	ok(t, err)

	_, err = redis.String(c.Do("SELECT", "5"))
	ok(t, err)

	_, err = redis.String(c.Do("SET", "foo", "baz"))
	ok(t, err)

	equals(t, "bar", s.Get("foo"))
	s.Select(5)
	equals(t, "baz", s.Get("foo"))

	// Another connection should have its own idea of the db:
	c2, err := redis.Dial("tcp", s.Addr())
	ok(t, err)
	v, err := redis.String(c2.Do("GET", "foo"))
	ok(t, err)
	equals(t, "bar", v)
}