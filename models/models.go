package models

// ExportRequest represents the request body for export endpoint
type ExportRequest struct {
	ScadContent string        `json:"scad_content" binding:"required" example:"cube([10,10,10]);"`
	Format      string        `json:"format" binding:"required" example:"png"`
	Options     ExportOptions `json:"options"`
}

// ExportOptions contains format-specific export options
type ExportOptions struct {
	PNG     *PNGOptions     `json:"png,omitempty"`
	STL     *STLOptions     `json:"stl,omitempty"`
	SVG     *SVGOptions     `json:"svg,omitempty"`
	PDF     *PDFOptions     `json:"pdf,omitempty"`
	ThreeMF *ThreeMFOptions `json:"3mf,omitempty"`
}

// PNGOptions contains PNG export options.
// Also used for webp and avif formats (which render via PNG internally).
type PNGOptions struct {
	Width  *int `json:"width,omitempty" example:"800"`
	Height *int `json:"height,omitempty" example:"600"`
}

// STLOptions contains STL export options
type STLOptions struct {
	DecimalPrecision *int `json:"decimal_precision,omitempty" example:"6" minimum:"1" maximum:"16"`
}

// SVGOptions contains SVG export options
type SVGOptions struct {
	Fill        *bool    `json:"fill,omitempty" example:"false"`
	FillColor   *string  `json:"fill_color,omitempty" example:"white"`
	Stroke      *bool    `json:"stroke,omitempty" example:"true"`
	StrokeColor *string  `json:"stroke_color,omitempty" example:"black"`
	StrokeWidth *float64 `json:"stroke_width,omitempty" example:"0.35"`
}

// PDFOptions contains PDF export options
type PDFOptions struct {
	PaperSize   *string  `json:"paper_size,omitempty" example:"a4" enums:"a6,a5,a4,a3,letter,legal,tabloid"`
	Orientation *string  `json:"orientation,omitempty" example:"portrait" enums:"portrait,landscape,auto"`
	ShowGrid    *bool    `json:"show_grid,omitempty" example:"false"`
	GridSize    *float64 `json:"grid_size,omitempty" example:"10"`
	Fill        *bool    `json:"fill,omitempty" example:"false"`
	FillColor   *string  `json:"fill_color,omitempty" example:"black"`
	Stroke      *bool    `json:"stroke,omitempty" example:"true"`
	StrokeColor *string  `json:"stroke_color,omitempty" example:"black"`
	StrokeWidth *float64 `json:"stroke_width,omitempty" example:"0.35"`
}

// ThreeMFOptions contains 3MF export options
type ThreeMFOptions struct {
	Unit              *string `json:"unit,omitempty" example:"millimeter" enums:"micron,millimeter,centimeter,meter,inch,foot"`
	DecimalPrecision  *int    `json:"decimal_precision,omitempty" example:"6" minimum:"1" maximum:"16"`
	Color             *string `json:"color,omitempty" example:"#f9d72c"`
	ColorMode         *string `json:"color_mode,omitempty" example:"model" enums:"model,none,selected-only"`
	MaterialType      *string `json:"material_type,omitempty" example:"color" enums:"color,basematerial"`
	AddMetadata       *bool   `json:"add_metadata,omitempty" example:"true"`
	MetadataTitle     *string `json:"metadata_title,omitempty" example:"My Model"`
	MetadataDesigner  *string `json:"metadata_designer,omitempty" example:"Designer Name"`
	MetadataDesc      *string `json:"metadata_description,omitempty" example:"Model description"`
	MetadataCopyright *string `json:"metadata_copyright,omitempty" example:"Copyright info"`
}

// SummaryRequest represents the request body for summary endpoint
type SummaryRequest struct {
	ScadContent string `json:"scad_content" binding:"required" example:"cube([10,10,10]);"`
	SummaryType string `json:"summary_type,omitempty" example:"all" enums:"all,cache,time,camera,geometry,bounding-box,area"`
}

// SummaryResponse represents the response from summary endpoint
type SummaryResponse struct {
	Summary map[string]interface{} `json:"summary"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"invalid parameter"`
	Message string `json:"message,omitempty" example:"detailed error message"`
}
