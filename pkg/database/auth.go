package database

// func UserRegister(ctx context.Context, input model.NewUser) (interface{}, error) {
// 	// Check Email
// 	_, err := FindUserByEmail(ctx, input.Email)
// 	if err == nil {
// 		// if err != record not found
// 		if err != gorm.ErrRecordNotFound {
// 			return nil, err
// 		}
// 	}

// 	_, err = UserCreate(ctx, input)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return nil, nil
// }

// func UserLogin(ctx context.Context, email string, password string) (interface{}, error) {
// 	getUser, err := FindUserByEmail(ctx, email)
// 	if err != nil {
// 		// if user not found
// 		if err == gorm.ErrRecordNotFound {
// 			return nil, &gqlerror.Error{
// 				Message: "Email not found",
// 			}
// 		}
// 		return nil, err
// 	}

// 	if err := getUser.ComparePassword(password); err != nil {
// 		return nil, err
// 	}

// 	return nil, nil
// }
