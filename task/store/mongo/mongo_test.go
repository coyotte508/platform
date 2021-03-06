package mongo_test

import (
	mgo "github.com/globalsign/mgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	logTest "github.com/tidepool-org/platform/log/test"
	storeStructuredMongo "github.com/tidepool-org/platform/store/structured/mongo"
	storeStructuredMongoTest "github.com/tidepool-org/platform/store/structured/mongo/test"
	taskStore "github.com/tidepool-org/platform/task/store"
	taskStoreMongo "github.com/tidepool-org/platform/task/store/mongo"
)

var _ = Describe("Mongo", func() {
	var cfg *storeStructuredMongo.Config
	var str *taskStoreMongo.Store
	var ssn taskStore.TaskSession

	BeforeEach(func() {
		cfg = storeStructuredMongoTest.NewConfig()
	})

	AfterEach(func() {
		if ssn != nil {
			ssn.Close()
		}
		if str != nil {
			str.Close()
		}
	})

	Context("New", func() {
		It("returns an error if unsuccessful", func() {
			var err error
			str, err = taskStoreMongo.NewStore(nil, nil)
			Expect(err).To(HaveOccurred())
			Expect(str).To(BeNil())
		})

		It("returns a new store and no error if successful", func() {
			var err error
			str, err = taskStoreMongo.NewStore(cfg, logTest.NewLogger())
			Expect(err).ToNot(HaveOccurred())
			Expect(str).ToNot(BeNil())
		})
	})

	Context("with a new store", func() {
		var mgoSession *mgo.Session

		BeforeEach(func() {
			var err error
			str, err = taskStoreMongo.NewStore(cfg, logTest.NewLogger())
			str.WaitUntilStarted()
			Expect(err).ToNot(HaveOccurred())
			Expect(str).ToNot(BeNil())
			mgoSession = storeStructuredMongoTest.Session().Copy()
		})

		AfterEach(func() {
			if mgoSession != nil {
				mgoSession.Close()
			}
		})

		Context("NewTaskSession", func() {
			It("returns a new session", func() {
				ssn = str.NewTaskSession()
				Expect(ssn).ToNot(BeNil())
			})
		})
	})
})
