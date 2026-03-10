package cases

import (
	"github.com/DuvanRozoParra/sicou/internal/case/handlers"
	"github.com/DuvanRozoParra/sicou/internal/case/repository"
	"github.com/DuvanRozoParra/sicou/internal/case/service"
	"github.com/DuvanRozoParra/sicou/pkg/database"
)

type Module struct {
	HandlerCaseFile     *handlers.SCaseFileHandler
	HandlerCaseEvent    *handlers.SCaseEventHandler
	handlerCheckListDef *handlers.SCheckListItemDef
	handlerServiceType  *handlers.SServiceTypeHandler
	handlerLegalArea    *handlers.SLegalAreaHandler
}

func NewModule() *Module {

	repo_caseFile := repository.NewSCaseFileRepository(database.DB)
	repo_caseEvent := repository.NewCaseEventRepository(database.DB)
	repo_check_list_def := repository.NewCheckListItemDefRepository(database.DB)
	repo_service_type := repository.NewServiceTypeRepository(database.DB)
	repo_legal_area := repository.NewLegalAreaRepository(database.DB)

	svc_caseFile := service.NewCaseFileService(repo_caseFile)
	svc_caseEvent := service.NewCaseEventService(repo_caseEvent)
	svc_check_list_def := service.NewCheckListItemDefService(repo_check_list_def)
	svc_service_type := service.NewServiceTypeService(repo_service_type)
	svc_legal_area := service.NewLegalAreaService(repo_legal_area)

	handler_caseFile := handlers.NewCaseFileHandler(svc_caseFile)
	handler_caseEvent := handlers.NewCaseEventHandler(svc_caseEvent)
	handler_check_list_def := handlers.NewCheckListItemDefHandler(svc_check_list_def)
	handler_service_type := handlers.NewServiceTypeHandler(svc_service_type)
	handler_legal_area := handlers.NewLegalAreaHandler(svc_legal_area)

	return &Module{
		HandlerCaseFile:     handler_caseFile,
		HandlerCaseEvent:    handler_caseEvent,
		handlerCheckListDef: handler_check_list_def,
		handlerServiceType:  handler_service_type,
		handlerLegalArea:    handler_legal_area,
	}
}
