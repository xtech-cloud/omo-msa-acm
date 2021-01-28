# omo-msa-acm
Micro Service Agent - acm

.PHONY: call-menu
call:
	MICRO_REGISTRY=consul micro call omo.msa.acm MenuService.AddOne '{"name":"api-1", "type":"g", "path":"omo/msa/acm/menu/add", "method":"get"}'
	MICRO_REGISTRY=consul micro call omo.msa.acm MenuService.AddOne '{"name":"api-2", "type":"g", "path":"omo/msa/acm/menu/add2", "method":"post"}'
	MICRO_REGISTRY=consul micro call omo.msa.acm MenuService.GetOne '{"uid":"5f0ec0949bf9162c1276a97f"}'
	MICRO_REGISTRY=consul micro call omo.msa.acm MenuService.RemoveOne '{"uid":"5f0ec0729bf9162c1276a97d" }'
	MICRO_REGISTRY=consul micro call omo.msa.acm MenuService.GetAll '{"uid":""}'
	MICRO_REGISTRY=consul micro call omo.msa.acm MenuService.UpdateBase '{"uid":"5f0ec0ac9bf9162c1276a980", "name":"api-5", "type":"hh", "path":"omo/msa/acm/menu/update", "method":""}'

.PHONY: call
call:
	MICRO_REGISTRY=consul micro call omo.msa.acm RoleService.AddOne '{"name":"checker26", "remark":"checker63 checker64", "menus":["5f0ec0949bf9162c1276a97f", "5f0ec0ac9bf9162c1276a980"]}'
	MICRO_REGISTRY=consul micro call omo.msa.acm RoleService.GetOne '{"uid":"5f0ec4d09bf9162c1276a98a"}'
	MICRO_REGISTRY=consul micro call omo.msa.acm RoleService.RemoveOne '{"uid":"5f0ec44e9bf9162c1276a988"}'
	MICRO_REGISTRY=consul micro call omo.msa.acm RoleService.GetAll '{"uid":""}'
	MICRO_REGISTRY=consul micro call omo.msa.acm RoleService.UpdateBase '{"uid":"5f0ec4d09bf9162c1276a98a", "name":"checker62", "remark":"test checker"}'
	MICRO_REGISTRY=consul micro call omo.msa.acm RoleService.AppendMenu '{"role":"5f0ec4d09bf9162c1276a98a", "menu":"5f0ec0ac9bf9162c1276a981"}'
	MICRO_REGISTRY=consul micro call omo.msa.acm RoleService.SubtractMenu '{"role":"5f0ec4d09bf9162c1276a98a", "menu":"5f0ec0ac9bf9162c1276a980"}'

.PHONY: call
call:
	MICRO_REGISTRY=consul micro call omo.msa.acm UserService.AddOne '{"user":"Tom18", "roles":["5f0ec4d09bf9162c1276a989","5f0ec4d09bf9162c1276a98a"]}'
	MICRO_REGISTRY=consul micro call omo.msa.acm UserService.GetOne '{"uid":"Tom4"}'
	MICRO_REGISTRY=consul micro call omo.msa.acm UserService.RemoveOne '{"uid":"Tom5"}'
	MICRO_REGISTRY=consul micro call omo.msa.acm UserService.GetList '{"page":2, "number":5}'
	MICRO_REGISTRY=consul micro call omo.msa.acm UserService.AppendRole '{"user":"Tom4", "role":"5f0ec5629bf9162c1276a98b"}'
	MICRO_REGISTRY=consul micro call omo.msa.acm UserService.SubtractRole '{"user":"Tom4", "role":"5f0ec4d09bf9162c1276a989"}'
	MICRO_REGISTRY=consul micro call omo.msa.acm UserService.IsPermission '{"user":"Tom4", "path":"omo/msa/acm/menu/add", "action":"get"}'