package repository

import (
	"context"
	"fmt"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockUserRepository struct {
    Users map[string]interface{} // Dùng map để lưu các user theo ID
    Err   error
}

// Find: Trả về tất cả người dùng phù hợp với filter (giả lập).
func (m *MockUserRepository) Find(ctx context.Context, filter bson.M) ([]models.UserDTO, error) {
    if m.Err != nil {
        return nil, m.Err
    }

    var results []models.UserDTO
    for _, user := range m.Users {
        // Kiểm tra và ép kiểu user (interface{}) thành models.UserDTO
        if u, ok := user.(models.UserDTO); ok {
            results = append(results, u)
        }
    }
    return results, nil
}

// FindOne: Tìm một người dùng dựa trên filter.
// func (m *MockUserRepository) FindOne(ctx context.Context, filter bson.M) *mongo.SingleResult {
//     userID, ok := filter["_id"].(primitive.ObjectID)
//     if !ok {
//         return mongo.NewSingleResultFromDocument(nil, fmt.Errorf("invalid filter format"), bson.NewRegistry())
//     }

//     user, found := m.Users[userID.Hex()]
//     registry := bson.NewRegistry()
//     if !found {
//         return mongo.NewSingleResultFromDocument(nil, mongo.ErrNoDocuments, registry)
//     }

//     return mongo.NewSingleResultFromDocument(user, nil, registry)
// }

// func (m *MockUserRepository) FindOne(ctx context.Context, filter bson.M) *mongo.SingleResult {
//     var userFound interface{}

//     for _, user := range m.Users {
//         userMap := user.(models.User) // Ép kiểu sang models.User

//         // Kiểm tra điều kiện tìm kiếm
//         match := true
//         for key, value := range filter {
//             switch key {
//             case "_id":
//                 if userMap.ID.Hex() != value.(primitive.ObjectID).Hex() {
//                     match = false
//                 }
//             case "username":
//                 if userMap.Username != value {
//                     match = false
//                 }
//             case "email":
//                 if userMap.Email != value {
//                     match = false
//                 }
//             }
//         }

//         if match {
//             userFound = user
//             break
//         }
//     }

//     registry := bson.NewRegistry()
//     if userFound == nil {
//         return mongo.NewSingleResultFromDocument(nil, mongo.ErrNoDocuments, registry)
//     }

//     return mongo.NewSingleResultFromDocument(userFound, nil, registry)
// }

func (m *MockUserRepository) FindOne(ctx context.Context, filter bson.M) *mongo.SingleResult {
    var userFound interface{}

    for _, user := range m.Users {
        userMap := user.(models.User) // Ép kiểu sang models.User

        // Kiểm tra điều kiện tìm kiếm
        match := true
        for key, value := range filter {
            switch key {
            case "_id":
                objID, ok := value.(primitive.ObjectID)
                if !ok {
                    match = false
                } else if userMap.ID.Hex() != objID.Hex() {
                    match = false
                }
            case "username":
                if userMap.Username != value.(string) {
                    match = false
                }
            case "email":
                if userMap.Email != value.(string) {
                    match = false
                }
            }
        }

        if match {
            userFound = user
            break
        }
    }

    registry := bson.NewRegistry()
    if userFound == nil {
        return mongo.NewSingleResultFromDocument(nil, mongo.ErrNoDocuments, registry)
    }

    return mongo.NewSingleResultFromDocument(userFound, nil, registry)
}


// DeleteOne: Xóa một người dùng dựa trên filter.
func (m *MockUserRepository) DeleteOne(ctx context.Context, filter bson.M) (*mongo.DeleteResult, error) {
    userID, ok := filter["_id"].(primitive.ObjectID)
    if !ok {
        return &mongo.DeleteResult{DeletedCount: 0}, fmt.Errorf("invalid filter format")
    }

    if _, found := m.Users[userID.Hex()]; found {
        delete(m.Users, userID.Hex())
        return &mongo.DeleteResult{DeletedCount: 1}, nil
    }

    return &mongo.DeleteResult{DeletedCount: 0}, nil
}
func (m *MockUserRepository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	userID, ok := filter["_id"].(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("invalid filter format")
	}

	userData, found := m.Users[userID.Hex()]
	if !found {
		return &mongo.UpdateResult{MatchedCount: 0}, nil
	}

	// Ép kiểu dữ liệu thành models.User
	user, ok := userData.(models.User)
	if !ok {
		return nil, fmt.Errorf("invalid user data type")
	}

	// Cập nhật từng trường dựa trên dữ liệu trong `update`
	if updates, ok := update["$set"].(bson.M); ok {
		for key, value := range updates {
			switch key {
			case "username":
				user.Username = value.(string)
			case "profile.full_name":
				user.Profile.FullName = value.(string)
			case "profile.phone_number":
				user.Profile.PhoneNumber = value.(string)
			case "profile.avatar_url":
				user.Profile.AvatarURL = value.(string)
			case "profile.bio":
				user.Profile.Bio = value.(string)
			case "profile.date_of_birth":
				user.Profile.DateOfBirth = value.(string)
            case "password":  // Cập nhật mật khẩu
				user.Password = value.(string)
            case "email":  
				user.Email = value.(string)
			case "is_active":  
				user.IsActive = value.(bool)
            }
		}
	}

	// Lưu lại dữ liệu đã cập nhật
	m.Users[userID.Hex()] = user

	return &mongo.UpdateResult{MatchedCount: 1}, nil
}


// ExistsByFilter: Kiểm tra xem có user nào khớp với filter không.
func (m *MockUserRepository) ExistsByFilter(ctx context.Context, filter bson.M) (bool, error) {
    if m.Err != nil {
        return false, m.Err
    }

    for _, user := range m.Users {
        // Giả lập logic kiểm tra filter, ví dụ lọc theo email hoặc username.
        userMap := user.(map[string]interface{})
        match := true
        for key, value := range filter {
            if userMap[key] != value {
                match = false
                break
            }
        }
        if match {
            return true, nil
        }
    }
    return false, nil
}

// InsertOne: Thêm một tài liệu mới vào mock repository.
func (m *MockUserRepository) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
    if m.Err != nil {
        return nil, m.Err
    }

    user, ok := document.(models.User)
    if !ok {
        return nil, fmt.Errorf("document is not of type models.User")
    }

    // Tạo một ObjectID giả lập nếu chưa có
    if user.ID.IsZero() {
        user.ID = primitive.NewObjectID()
    }

    // Lưu tài liệu vào map
    m.Users[user.ID.Hex()] = user

    return &mongo.InsertOneResult{InsertedID: user.ID}, nil
}
