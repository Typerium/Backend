package handlers

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	gql "github.com/99designs/gqlgen/graphql"
	"github.com/gramework/gramework"
	"github.com/pkg/errors"

	"typerium/internal/app/gateway/graphql"
)

type graphqlDecodeFunc func() (*gramework.GQLRequest, error)

func Handler(executor graphql.Executor) gramework.RequestHandler {
	return func(ctx *gramework.Context) {
		contentType := strings.Split(strings.ToLower(ctx.ContentType()), ";")
		if len(contentType) == 0 {
			ctx.BadRequest()
			return
		}

		var decodeFunc graphqlDecodeFunc
		switch contentType[0] {
		case "application/json", "application/graphql":
			decodeFunc = ctx.DecodeGQL
		case "multipart/form-data":
			decodeFunc = multipartDecodeFunc(ctx)
		default:
			return
		}

		req, err := decodeFunc()
		if err != nil {
			if err == gramework.ErrInvalidGQLRequest {
				ctx.BadRequest(err)
				return
			}
			ctx.Logger.WithError(err).Error("can't decode request")
			ctx.Err500()
			return
		}

		gqlCtx := context.Background()

		resp := executor.Exec(gqlCtx, &graphql.Request{
			Query:         req.Query,
			OperationName: req.OperationName,
			Variables:     req.Variables,
		})

		err = ctx.JSON(resp)
		if err != nil {
			ctx.Err500()
		}
	}
}

const (
	variablesFilesKey = "variables"
)

func multipartDecodeFunc(ctx *gramework.Context) graphqlDecodeFunc {
	return func() (req *gramework.GQLRequest, err error) {
		data, err := ctx.Request.MultipartForm()
		if err != nil {
			err = errors.WithStack(err)
			return
		}

		operations, ok := data.Value["operations"]
		if !ok || len(operations) != 1 {
			return nil, errors.WithStack(gramework.ErrInvalidGQLRequest)
		}

		req = new(gramework.GQLRequest)
		err = json.NewDecoder(strings.NewReader(operations[0])).Decode(req)
		if err != nil {
			return nil, errors.WithStack(gramework.ErrInvalidGQLRequest)
		}

		filesMap, ok := data.Value["map"]
		if !ok || len(operations) != 1 {
			return
		}

		files := make(map[string][]string)
		err = json.NewDecoder(strings.NewReader(filesMap[0])).Decode(&files)
		if err != nil {
			return nil, errors.WithStack(gramework.ErrInvalidGQLRequest)
		}

		for key, value := range files {
			if len(value) != 1 {
				continue
			}

			var variableName string
			indexVariable := -1
			path := strings.Split(value[0], ".")
			switch len(path) {
			case 2:
				if path[0] != variablesFilesKey {
					break
				}
				variableName = path[1]
			case 3:
				if path[0] == variablesFilesKey {
					variableName = path[1]
					indexVariable, err = strconv.Atoi(path[2])
					if err != nil {
						continue
					}
					break
				}
				if path[1] == variableName {
					variableName = path[2]
				}
			case 4:
				if path[1] != variableName {
					break
				}
				variableName = path[2]
				indexVariable, err = strconv.Atoi(path[3])
				if err != nil {
					continue
				}
			}
			if len(variableName) == 0 {
				continue
			}

			files, ok := data.File[key]
			if !ok || len(files) != 1 {
				return nil, errors.WithStack(gramework.ErrInvalidGQLRequest)
			}
			fileParams := files[0]
			file, err := fileParams.Open()
			if err != nil {
				return nil, errors.WithStack(err)
			}
			upload := gql.Upload{
				Filename: fileParams.Filename,
				Size:     fileParams.Size,
				File:     file,
			}
			variableValue := req.Variables[variableName]

			variableValueArr, ok := variableValue.([]interface{})
			if !ok {
				if indexVariable != -1 {
					return nil, errors.WithStack(gramework.ErrInvalidGQLRequest)
				}

				variableValue = upload
			}

			if indexVariable == -1 {
				indexVariable = 0
			}
			if indexVariable >= len(variableValueArr) {
				return nil, errors.WithStack(gramework.ErrInvalidGQLRequest)
			}
			variableValueArr[indexVariable] = upload
		}

		return
	}
}
