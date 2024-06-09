with

source as (

    select * from {{ "{{" }} source('{{.Schema}}', '{{.Name}}') {{ "}}" }}

),

renamed as (

    select
{{- $j := 0 -}}
{{- range $group, $columns := .DataTypeGroups -}}
    {{- if eq $group "text" }}
        -- text
        {{ range $i, $e := $columns -}}
            {{- if or (ne $j 0) (ne $i 0) -}}, {{- end -}}
            {{- .Name }} as {{ .Name | lower }}
        {{ end -}}
    {{- end -}}
    {{- if eq $group "numbers" }}
        -- numbers
        {{ range $i, $e := $columns -}}
            {{- if or (ne $j 0) (ne $i 0) -}}, {{- end -}}
            {{- .Name }} as {{ .Name | lower }}
        {{ end -}}
    {{- end -}}
    {{- if eq $group "booleans" }}
        -- booleans
        {{ range $i, $e := $columns -}}
            {{- if or (ne $j 0) (ne $i 0) -}}, {{- end -}}
            {{- .Name }} as {{ .Name | lower }}
        {{ end -}}
    {{- end -}}
    {{- if eq $group "datetimes" }}
        -- datetimes
        {{ range $i, $e := $columns -}}
            {{- if or (ne $j 0) (ne $i 0) -}}, {{- end -}}
            {{- .Name }} as {{ .Name | lower }}
        {{ end -}}
    {{- end -}}
    {{- if eq $group "timestamps" }}
        -- timestamps
        {{ range $i, $e := $columns -}}
            {{- if or (ne $j 0) (ne $i 0) -}}, {{- end -}}
            {{- .Name }} as {{ .Name | lower }}
        {{ end -}}
    {{- end -}}
{{- $j = 1 -}}
{{ end }}
    from source
)

select * from renamed
