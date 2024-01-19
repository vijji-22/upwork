package rbac

import (
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func encryptPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func comparePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// convertToAccessWithReferenceIDMap will convert AccessWithReferenceID slice to AccessWithReferenceIDMap
// it will merge all projects by reference id
// if reference id is nil then it will merge all projects in GlobalProject
// in AccessWithReferenceIDMap GlobalProject will have more priority than Reference when we read project
func convertToAccessWithReferenceIDMap[IDTYPE int64 | string](accessWithReferenceID []AccessWithReferenceID[IDTYPE]) *AccessWithReferenceIDMap[IDTYPE] {
	var accessWithReferenceIDMap AccessWithReferenceIDMap[IDTYPE]

	accessWithReferenceIDMap.Reference = make(map[IDTYPE]*utils.Set[string])
	for _, access := range accessWithReferenceID {
		accessWithReferenceIDMap.Name = access.AccessName
		accessWithReferenceIDMap.ReferenceKey = access.ReferenceKey

		if access.ReferenceID == nil {
			if accessWithReferenceIDMap.GlobalProject == nil {
				accessWithReferenceIDMap.GlobalProject = utils.NewSet[string]()
			}
			accessWithReferenceIDMap.GlobalProject.Add(access.Project.Slice()...)
			continue
		}
		if _, ok := accessWithReferenceIDMap.Reference[*access.ReferenceID]; !ok {
			accessWithReferenceIDMap.Reference[*access.ReferenceID] = utils.NewSet[string]()
		}

		accessWithReferenceIDMap.Reference[*access.ReferenceID].Add(access.Project.Slice()...)
	}

	return &accessWithReferenceIDMap
}
