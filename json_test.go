package gonfig_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/payneio/gonfig"
)

var _ = Describe("JsonConfig", func() {
	var (
		err error
		cfg WritableConfig
	)

	BeforeEach(func() {
		cfg = NewJsonConfig("./config_valid.json")
		err = cfg.Load()
	})

	Context("When the JSON config marshals properly", func() {
		It("Should have the variables in config", func() {
			Expect(cfg.Get("test")).To(Equal("123"))
		})
		It("Should have the nested variables in config", func() {
			fmt.Println("sup?", cfg.All())
			Expect(cfg.Get("test_object:test_b")).To(Equal("b"))
		})
		It("Should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("When config fails to marshal", func() {
		BeforeEach(func() {
			cfg = NewJsonConfig("./config_invalid.json")
			err = cfg.Load()
		})
		It("should return a functional config", func() {
			Expect(cfg).ToNot(BeZero())
			cfg.Set("QQ", "123")
			Expect(cfg.Get("QQ")).To(Equal("123"))
		})

		It("should error", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Context("When the JSON config does not exist", func() {
		BeforeEach(func() {
			cfg = NewJsonConfig("./config_nonexisting.json")
			err = cfg.Load()
		})
		It("should return a functional config", func() {
			Expect(cfg).ToNot(BeZero())
			cfg.Set("QQ", "123")
			Expect(cfg.Get("QQ")).To(Equal("123"))
		})
		It("should error", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Config conversion", func() {
		It("Should be possible ro construct new JSON config from a gonfig hierarchy", func() {
			cfg := NewConfig(nil)
			cfg.Use("config_a", NewMemoryConfig())
			cfg.Use("config_b", NewMemoryConfig())
			cfg.Use("config_a").Set("config_a_var_a", "conf_a")
			cfg.Use("config_b").Set("config_b_var_a", "conf_b")
			jsonConf := &JsonConfig{cfg, "./config.json"}
			err := jsonConf.Save()
			Expect(err).ToNot(HaveOccurred())
			jsonConf2 := NewJsonConfig("./config.json", cfg)
			err = jsonConf2.Save()
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Namespacing", func() {

	})
})
