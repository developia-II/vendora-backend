package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}
type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Email            string             `json:"email" bson:"email" validate:"required,email"`
	Name             string             `json:"name" bson:"name"`
	Address          string             `json:"address" bson:"address"`
	Phone            string             `json:"phone" bson:"phone"`
	Password         string             `json:"-" bson:"password" validate:"required,min=6"`
	Role             string             `json:"role" bson:"role"`
	CreatedAt        time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time          `json:"updatedAt" bson:"updatedAt"`
	IsVerified       bool               `json:"isverified" bson:"isverified"`
	ResetToken       string             `json:"-" bson:"resetToken,omitempty"`
	ResetTokenExpiry time.Time          `json:"-" bson:"resetTokenExpiry,omitempty"`
	PasswordResetAt  time.Time          `json:"-" bson:"passwordResetAt,omitempty"`

	OnboardingCompleted bool             `json:"onboardingCompleted" bson:"onboardingCompleted"`
	Preferences         *UserPreferences `json:"preferences" bson:"preferences"`
	Interests           *UserInterests   `json:"interests" bson:"interests"`
	Profile             *UserProfile     `json:"profile" bson:"profile"`
	VendorStatus        string           `json:"vendorStatus" bson:"vendorStatus"` // "", "pending", "approved", "rejected"

	SellerApplication *SellerApplication `json:"sellerApplication" bson:"sellerApplication"`
}

type UserPreferences struct {
	Categories        []string         `json:"categories,omitempty" bson:"categories,omitempty"`
	BudgetRange       string           `json:"budgetRange" bson:"budgetRange"`
	ShoppingFrequency string           `json:"shoppingFrequency" bson:"shoppingFrequency"`
	SpecialPrefs      *map[string]bool `json:"specialPrefs,omitempty" bson:"specialPrefs,omitempty"`
}

type UserInterests struct {
	Categories []string `json:"categories" bson:"categories" validate:"required,min=1,max=3,dive,required"`
	IsSet      bool     `json:"isSet" bson:"isSet"`
}

type UserProfile struct {
	Location       string `json:"location" bson:"location"`
	Bio            string `json:"bio" bson:"bio"`
	ProfilePicture string `json:"profileImage,omitempty" bson:"profileImage,omitempty"`
}

type SellerApplication struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty"`
	UserID       primitive.ObjectID  `bson:"userID"`
	BusinessName string              `json:"businessName" bson:"businessName"`
	BusinessType string              `json:"businessType" bson:"businessType"`
	Categories   []string            `json:"categories" bson:"categories"`
	Description  string              `json:"description" bson:"description"`
	Location     string              `json:"location" bson:"location"`
	Status       string              `json:"status" bson:"status"` // pending, approved, rejected
	AppliedAt    time.Time           `json:"appliedAt" bson:"appliedAt"`
	ReviewedAt   *time.Time          `json:"reviewedAt" bson:"reviewedAt"`
	ReviewedBy   *primitive.ObjectID `json:"reviewedBy" bson:"reviewedBy"`
}

type VendorApplication struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty"`
	UserID              primitive.ObjectID `bson:"userID"`
	BusinessName        string             `json:"businessName" bson:"businessName"`
	BusinessType        string             `json:"businessType" bson:"businessType"`
	BusinessDescription string             `json:"businessDescription" bson:"businessDescription"`
	ContactEmail        string             `json:"contactEmail" bson:"contactEmail"`
	ContactPhone        string             `json:"contactPhone" bson:"contactPhone"`
	BusinessAddress     string             `json:"businessAddress" bson:"businessAddress"`
	TaxID               string             `json:"taxId" bson:"taxId"`
	Website             string             `json:"website" bson:"website"`
	SocialMedia         []string           `json:"socialMedia" bson:"socialMedia"`
	Products            []string           `json:"products" bson:"products"`
	Experience          string             `json:"experience" bson:"experience"`
	Motivation          string             `json:"motivation" bson:"motivation"`
	Status              string             `json:"status" bson:"status"` // pending, approved, rejected
	AppliedAt           time.Time          `json:"appliedAt" bson:"appliedAt"`
	ReviewedAt          *time.Time         `json:"reviewedAt" bson:"reviewedAt"`
	ReviewedBy          string             `json:"reviewedBy" bson:"reviewedBy"`
	ReviewNotes         string             `json:"reviewNotes" bson:"reviewNotes"`
}
type UserOnboardingDraft struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty"`
	UserID        primitive.ObjectID     `json:"userID" bson:"userID"`
	Role          string                 `json:"role" bson:"role"`
	Step          int                    `json:"step" bson:"step"`
	StepCompleted bool                   `json:"stepCompleted" bson:"stepCompleted"`
	StepData      map[string]interface{} `json:"stepData" bson:"stepData"`
	UpdatedAt     time.Time              `json:"updatedAt" bson:"updatedAt"`
	Version       int                    `json:"version" bson:"version"`
}
