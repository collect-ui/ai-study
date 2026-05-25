select a.*
from sys_code a
left join sys_code parent_code
  on parent_code.sys_code_type = replace(a.sys_code_type, 'knowledge_category_', 'knowledge_type_')
  and parent_code.sys_code = ifnull(a.icon, '')
where a.sys_code_type = {{.sys_code_type}}
{{ if .sys_code_list }}
and a.sys_code in ({{.sys_code_list}})
{{ end }}
order by ifnull(parent_code.order_index, a.order_index), a.order_index, a.sys_code
