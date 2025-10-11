package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/developia-II/ecommerce-backend/internal/models"
	"github.com/developia-II/ecommerce-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OnboardingHandler struct {
	DB *mongo.Database
}

func NewOnboardingHandler(db *mongo.Database) *OnboardingHandler {
	return &OnboardingHandler{DB: db}
}

var onboardingValidator = validator.New()

func (h *OnboardingHandler) ClientUpdateInterest(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Missing or invalid token"))
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Missing or Invalid token"))
		return
	}
	userId := claims.UserID
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	var user models.User
	collection := h.DB.Collection("users")
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Invalid or missing token"))
		return
	}
	filter := bson.M{"_id": objectId}
	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("No user found"))
		return
	}

	var userInterest models.UserInterests
	if err := c.ShouldBindJSON(&userInterest); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid JSON format"))
		return
	}
	if err := onboardingValidator.Struct(userInterest); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Validation failed: "+err.Error()))
		return
	}

	// Initialize interests object if it's null, then set the fields
	update := bson.M{
		"$set": bson.M{
			"interests": bson.M{
				"categories": userInterest.Categories,
				"isSet":      true,
			},
			"updatedAt": time.Now(),
		},
	}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update user interest: "+err.Error()))
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Interests updated successfully", gin.H{
		"categories":    userInterest.Categories,
		"isInterestSet": true,
	}))
}

func (h *OnboardingHandler) ClientUpdatePreference(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid or missing token"))
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.VerifyToken(tokenString)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to verify token"))
		return
	}

	objectId, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid token"))
		return
	}

	var userPref models.UserPreferences

	if err := c.ShouldBindJSON(&userPref); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Validation failed: "+err.Error()))
		return
	}
	if err := onboardingValidator.Struct(&userPref); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Validation failed: "+err.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	collection := h.DB.Collection("users")

	var user models.User
	if err := collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("No user found"))
		return
	}
	filter := bson.M{"_id": objectId}

	// Initialize preferences object if it's null, then set the fields
	update := bson.M{
		"$set": bson.M{
			"preferences": bson.M{
				"budgetRange":       userPref.BudgetRange,
				"shoppingFrequency": userPref.ShoppingFrequency,
				"specialPrefs":      userPref.SpecialPrefs,
			},
			"updatedAt": time.Now(),
		},
	}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update user: "+err.Error()))
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("User preference updated successfully", gin.H{
		"budgetRange":       userPref.BudgetRange,
		"shoppingFrequency": userPref.ShoppingFrequency,
		"specialPrefs":      userPref.SpecialPrefs,
	}))

}

func (h *OnboardingHandler) CompleteOnboardingFlow(c *gin.Context) {

	authHeader := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid or missing token"))
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid or expired token"))
		return
	}

	location := c.PostForm("location")
	bio := c.PostForm("bio")
	file, err := c.FormFile("profile_picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Profile picture is required"))
		return
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Only JPEG/PNG images are allowed"))
		return
	}

	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to initialize Cloudinary"))
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to process image"))
		return
	}
	defer src.Close()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	uploadResult, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder: "users/profiles",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to upload image"))
		return
	}

	objectId, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid user ID"))
		return
	}

	collection := h.DB.Collection("users")
	filter := bson.M{"_id": objectId}

	// Initialize profile object if it's null, then set the fields
	update := bson.M{
		"$set": bson.M{
			"profile": bson.M{
				"location":     location,
				"bio":          bio,
				"profileImage": uploadResult.SecureURL,
			},
			"onboardingCompleted": true,
			"updatedAt":           time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update profile: "+err.Error()))
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Onboarding completed successfully", gin.H{
		"success":             true,
		"onboardingCompleted": true,
		"profile": gin.H{
			"location":     location,
			"bio":          bio,
			"profileImage": uploadResult.SecureURL,
		},
	}))
}

func (h *OnboardingHandler) UserOnboardingDraft(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid or missing token"))
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Failed to verify token"))
		return
	}
	userId := claims.UserID

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid or missing token"))
		return
	}
	var input models.UserOnboardingDraft
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
		return
	}
	if err := onboardingValidator.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Validation failed: "+err.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	collection := h.DB.Collection("drafts")

	filter := bson.M{"userID": objectId, "role": input.Role}
	update := bson.M{
		"$set": bson.M{
			"step":          input.Step,
			"stepCompleted": input.StepCompleted,
			"stepData":      input.StepData,
			"role":          input.Role,
			"updatedAt":     time.Now(),
		},
		"$setOnInsert": bson.M{
			"userID": objectId,
		},
		"$inc": bson.M{"version": 1},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	_, _ = h.DB.Collection("drafts").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "userID", Value: 1}, {Key: "role", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	var saved models.UserOnboardingDraft
	if err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&saved); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to save draft"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Draft saved successfully", gin.H{
		"success": true,
		"data": gin.H{
			"id":            saved.ID,
			"userID":        saved.UserID,
			"role":          saved.Role,
			"step":          saved.Step,
			"stepCompleted": saved.StepCompleted,
			"stepData":      saved.StepData,
			"version":       saved.Version,
			"updatedAt":     saved.UpdatedAt,
		},
	}))
}

func (h *OnboardingHandler) GetOnboardingDraft(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid or missing token"))
		return
	}
	claims, err := utils.VerifyToken(strings.TrimPrefix(authHeader, "Bearer "))
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid or expired token"))
		return
	}
	role := c.Query("role")
	if role == "" {
		role = "customer"
	}

	objectId, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid user ID"))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var draft models.UserOnboardingDraft
	err = h.DB.Collection("drafts").FindOne(ctx, bson.M{
		"userID": objectId,
		"role":   role,
	}).Decode(&draft)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusOK, utils.SuccessResponse("No draft found", gin.H{"success": true, "data": nil}))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch draft"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Draft fetched successfully", gin.H{
		"success": true,
		"data": gin.H{
			"id":            draft.ID,
			"userID":        draft.UserID,
			"role":          draft.Role,
			"step":          draft.Step,
			"stepCompleted": draft.StepCompleted,
			"stepData":      draft.StepData,
			"version":       draft.Version,
			"updatedAt":     draft.UpdatedAt,
		},
	}))
}

func (h *OnboardingHandler) SellerBusinessType(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid or missing token"))
		return
	}
	claims, err := utils.VerifyToken(strings.TrimPrefix(auth, "Bearer "))
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid or missing token"))
		return
	}

	var input models.SellerBusinessInfo

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid json format"))
		return
	}
	if err := onboardingValidator.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid json payload: "+err.Error()))
		return
	}

	userID, _ := primitive.ObjectIDFromHex(claims.UserID)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Save to drafts (role=vendor)
	filter := bson.M{"userID": userID, "role": "vendor"}
	update := bson.M{
		"$set": bson.M{
			"stepData.businessInfo": input,
			"updatedAt":             time.Now(),
		},
		"$setOnInsert": bson.M{"userID": userID, "role": "vendor", "step": 1},
		"$inc":         bson.M{"version": 1},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var saved models.UserOnboardingDraft
	if err := h.DB.Collection("drafts").FindOneAndUpdate(ctx, filter, update, opts).Decode(&saved); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to save: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Business details updated", gin.H{
		"success": true,
		"data":    gin.H{"businessInfo": input},
	}))
}

func (h *OnboardingHandler) SellerBusinessCategory(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid or missing token"))
		return
	}
	claims, err := utils.VerifyToken(strings.TrimPrefix(auth, "Bearer "))
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid or missing token"))
		return
	}

	type categoryInput struct {
		Categories []string `json:"categories" validate:"required,min=1,max=5,dive,required"`
	}
	var input categoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid json format"))
		return
	}
	if err := onboardingValidator.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid json payload: "+err.Error()))
		return
	}

	userID, _ := primitive.ObjectIDFromHex(claims.UserID)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	filter := bson.M{"userID": userID, "role": "vendor"}
	update := bson.M{
		"$set": bson.M{
			"stepData.categories": input.Categories,
			"updatedAt":           time.Now(),
		},
		"$setOnInsert": bson.M{"userID": userID, "role": "vendor", "step": 2},
		"$inc":         bson.M{"version": 1},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var saved models.UserOnboardingDraft
	if err := h.DB.Collection("drafts").FindOneAndUpdate(ctx, filter, update, opts).Decode(&saved); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to save: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Business categories updated", gin.H{
		"success": true,
		"data":    gin.H{"categories": input.Categories},
	}))
}

func (h *OnboardingHandler) SellerBusinessInfo(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid or missing token"))
		return
	}
	claims, err := utils.VerifyToken(strings.TrimPrefix(auth, "Bearer "))
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid or missing token"))
		return
	}
	var input models.BusinessDetails
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid json payload"))
		fmt.Println("Error", err)
		return
	}
	if err := onboardingValidator.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid json payload"))
		fmt.Println("Error", err)
		return
	}
	userID, _ := primitive.ObjectIDFromHex(claims.UserID)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"userID": userID,
		"role":   "vendor",
	}
	update := bson.M{
		"$set": bson.M{
			"stepData.businessDetails": bson.M{
				"businessName": input.BusinessName,
				"description":  input.Description,
				"location":     input.Location,
				"url":          input.Url,
			},
			"updatedAt": time.Now(),
		},
		"$setOnInsert": bson.M{"userID": userID, "role": "vendor", "step": 3},
		"$inc":         bson.M{"version": 1},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var saved models.UserOnboardingDraft
	if err := h.DB.Collection("drafts").FindOneAndUpdate(ctx, filter, update, opts).Decode(&saved); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to save: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Business details updated", gin.H{
		"success": true,
		"data": gin.H{
			"businessName": input.BusinessName,
			"description":  input.Description,
			"location":     input.Location,
			"url":          input.Url,
		},
	}))
}

func (h *OnboardingHandler) StoreDetails(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid or missing token"))
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Invalid or expired token"))
		return
	}
	userID, _ := primitive.ObjectIDFromHex(claims.UserID)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	storeName := c.PostForm("storeName")
	storeDescription := c.PostForm("storeDescription")
	primaryColor := c.PostForm("primaryColor")
	accentColor := c.PostForm("accentColor")
	file, err := c.FormFile("storeLogo")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid file"))
		return
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Only JPEG/PNG images are allowed"))
		return
	}
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to initialize Cloudinary"))
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to process image"))
		return
	}
	defer src.Close()

	uploadResult, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder: "stores/logo",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to upload image"))
		return
	}

	collection := h.DB.Collection("drafts")
	filter := bson.M{
		"userID": userID,
		"role":   "vendor",
	}
	update := bson.M{
		"$set": bson.M{
			"stepData.storeDetails": bson.M{
				"storeLogo":        uploadResult.SecureURL,
				"storeName":        storeName,
				"storeDescription": storeDescription,
				"primaryColor":     primaryColor,
				"accentColor":      accentColor,
			},
			"updatedAt": time.Now(),
		},
		"$setOnInsert": bson.M{"userID": userID, "role": "vendor", "step": 4},
		"$inc":         bson.M{"version": 1},
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var saved models.UserOnboardingDraft
	if err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&saved); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update store details"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Store details updated", gin.H{
		"success": true,
		"data": gin.H{
			"storeLogo":        uploadResult.SecureURL,
			"storeName":        storeName,
			"storeDescription": storeDescription,
			"primaryColor":     primaryColor,
			"accentColor":      accentColor,
		},
	}))
}
