package cache

//Mock to mock cache database
//type Mock struct {
//	mock.Mock
//}
//
////Initialize to run before tests
//func (rc *Mock) Initialize(map[string]string) {
//	instance = rc
//}
//
////Get to mock Get calls
//func (rc *Mock) Get(key string, object interface{}) error {
//	args := rc.Called(key, object)
//	return args.Error(0)
//}
//
////Set to mock Set calls
//func (rc *Mock) Set(k string, obj interface{}, d time.Duration) error {
//	args := rc.Called(k, obj, d)
//	return args.Error(0)
//}
//
////Delete to mock Delete calls
//func (rc *Mock) Delete(key string) error {
//	args := rc.Called(key)
//	return args.Error(0)
//}
//
////Disconnect does nothing
//func (rc *Mock) Disconnect() {}
