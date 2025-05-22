package model


type FindAllEventCriteria struct {
    Limit  int   `json:"limit" validate:"gte=0"`
    Offset int   `json:"offset" validate:"gte=0"`
    Search string `json:"search"`
}
