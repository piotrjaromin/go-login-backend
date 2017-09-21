package accounts

import (
        "github.com/piotrjaromin/go-login-backend/web"
        "github.com/piotrjaromin/go-login-backend/dal"
        "time"
        "github.com/satori/go.uuid"
        "github.com/op/go-logging"
)

type Dal struct {
        GetById                func(id string) (PasswordlessAccount, error)
        GetByEmail             func(email string) (PasswordlessAccount, error)
        GetWithPasswordById    func(id string) (SecuredAccount, error)
        GetWithPasswordByEmail func(email string) (SecuredAccount, error)
        getByFbID              func(fbId string) (PasswordlessAccount, error)
        UpdateByEmail          func(email string, handleUpdateFunc func(*SecuredAccount) error) (error)
        updateByID             func(id string, handleUpdateFunc func(*SecuredAccount) error) (error)
        CreateAccount          func(secAccount SecuredAccount) (string, error)
}

func CreateDal(accountsRepo dal.Dal) Dal {

        var log = logging.MustGetLogger("[AccountDal]")
        getWithPasswordById := func(id string) (SecuredAccount, error) {

                log.Debug("Getting account with password field")
                account := new(SecuredAccount)
                error := accountsRepo.GetById(id, account)

                if len(account.Id) == 0 {
                        return *account, ErrAccountNotFound
                }

                return *account, error
        }

        getById := func(id string) (PasswordlessAccount, error) {
                acc, err := getWithPasswordById(id)
                return acc.PasswordlessAccount, err
        }

        getByQuery := func(query dal.Query, pagination web.Pagination) ([]SecuredAccount, error) {
                accs := make([]SecuredAccount, 0)
                if err := accountsRepo.GetByQuery(&accs, pagination, query); err != nil {
                        return []SecuredAccount{}, err
                }

                return accs, nil
        }

        getWithPasswordByEmail := func(email string) (SecuredAccount, error) {

                query := dal.NewQueryBuilder().WithField("email", email).Build()
                accs, err := getByQuery(query, web.DefaultPagination())
                if err != nil {
                        return SecuredAccount{}, err
                }

                if len(accs) == 0 {
                        return SecuredAccount{}, ErrAccountNotFound
                }

                return accs[0], nil
        }

        getByFbID := func(fbId string) (PasswordlessAccount, error) {

                query := dal.NewQueryBuilder().WithField("providers.fb", fbId).Build()
                accs, err := getByQuery(query, web.DefaultPagination())

                if err != nil {
                        return PasswordlessAccount{}, err
                }

                if len(accs) == 0 {
                        return PasswordlessAccount{}, ErrAccountNotFound
                }

                return accs[0].PasswordlessAccount, nil
        }

        getByEmail := func(email string) (PasswordlessAccount, error) {
                acc, err := getWithPasswordByEmail(email)
                return acc.PasswordlessAccount, err
        }

        updateByEmail := func(email string, updateHandle func(*SecuredAccount) error) error {

                acc, err := getWithPasswordByEmail(email)
                if err != nil {
                        return err
                }

                updateHandle(&acc)
                return accountsRepo.UpdateByQuery(dal.NewQueryBuilder().WithField("email", email).Build(), acc)
        }

        updateByID := func(id string, updateHandle func(*SecuredAccount) error) error {

                acc, err := getWithPasswordById(id)
                if err != nil {
                        return err
                }

                updateHandle(&acc)
                return accountsRepo.Update(id, acc)
        }

        createAccount := func(secAccount SecuredAccount) (string, error) {

                secAccount.CreatedAt = time.Now()
                secAccount.Id = uuid.NewV4().String()

                _, saveAccErr := accountsRepo.Save(secAccount)
                return secAccount.Id, saveAccErr

        }

        return Dal{
                GetById: getById,
                GetByEmail: getByEmail,
                GetWithPasswordById: getWithPasswordById,
                GetWithPasswordByEmail: getWithPasswordByEmail,
                UpdateByEmail: updateByEmail,
                updateByID: updateByID,
                getByFbID: getByFbID,
                CreateAccount: createAccount,
        }

}
