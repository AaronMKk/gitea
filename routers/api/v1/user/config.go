package user

import (
	"github.com/opensourceways/xihe-server/utils"
	"github.com/sirupsen/logrus"
)

var (
	apiConfig         APIConfig
	encryptHelper     utils.SymmetricEncryption
	encryptHelperCSRF utils.SymmetricEncryption
	private_log       *logrus.Entry
)

type Tags struct {
	ModelTagDomains         []string `json:"model"            required:"true"`
	ProjectTagDomains       []string `json:"project"          required:"true"`
	DatasetTagDomains       []string `json:"dataset"          required:"true"`
	GlobalModelTagDomains   []string `json:"global_model"     required:"true"`
	GlobalProjectTagDomains []string `json:"global_project"   required:"true"`
	GlobalDatasetTagDomains []string `json:"global_dataset"   required:"true"`
}

type APIConfig struct {
	Tags                           Tags   `json:"tags"                        required:"true"`
	TokenKey                       string `json:"token_key"                   required:"true"`
	TokenExpiry                    int64  `json:"token_expiry"                required:"true"`
	EncryptionKey                  string `json:"encryption_key"              required:"true"`
	EncryptionKeyForCSRF           string `json:"encryption_key_csrf"         required:"true"`
	EncryptionKeyForGitlabToken    string `json:"encryption_key_gitlab_token" required:"true"`
	DefaultPassword                string `json:"default_password"            required:"true"`
	MaxTrainingRecordNum           int    `json:"max_training_record_num"     required:"true"`
	InferenceDir                   string `json:"inference_dir"`
	InferenceBootFile              string `json:"inference_boot_file"`
	InferenceTimeout               int    `json:"inference_timeout"`
	PodTimeout                     int    `json:"pod_timeout"`
	MaxPictureSizeToDescribe       int64  `json:"-"`
	MaxPictureSizeToVQA            int64  `json:"-"`
	MaxCompetitionSubmmitFileSzie  int64  `json:"max_competition_submmit_file_size"`
	MinSurvivalTimeOfInference     int    `json:"min_survival_time_of_inference"`
	MaxTagsNumToSearchResource     int    `json:"max_tags_num_to_search_resource"`
	MaxTagKindsNumToSearchResource int    `json:"max_tag_kinds_num_to_search_resource"`
}
