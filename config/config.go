package config

// Configuration for app
type Configuration struct {
	CatID                          string `mapstructure:"CAT_ID"`
	CatName                        string `mapstructure:"CAT_NAME"`
	FirebaseCredentials            string `mapstructure:"GOOGLE_FIREBASE_CREDENTIAL_FILE"`
	FirestoreCollectionSource      string `mapstructure:"GOOGLE_FIRESTORE_COLLECTION_SOURCE"`
	FirestoreCollectionDestination string `mapstructure:"GOOGLE_FIRESTORE_COLLECTION_DESTINATION"`
}
