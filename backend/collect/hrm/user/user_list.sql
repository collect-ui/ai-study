select
    a.user_id,
    a.nick,
    a.username,
    {{ if .with_password }}
    a.password as password,
    {{ else }}
    '' as password,
    {{ end }}
    a.status,
    a.entry_date,
    a.email,
    a.phone,
    a.leave_date,
    a.is_delete,
    a.create_time,
    a.create_user,
    a.modify_user,
    a.modify_time,
    a.work_code,
    a.leave_reason,
    a.sex,
    (
        select group_concat(ur.role_name)
        from role ur
        left join user_role_id_list r on r.role_id = ur.role_id
        where r.user_id = a.user_id
    ) as role_names,
    (
        select group_concat(ur.role_code)
        from role ur
        left join user_role_id_list r on r.role_id = ur.role_id
        where r.user_id = a.user_id
    ) as roles,
    c.sys_code_text as status_text,
    sex.sys_code_text as sex_name
from user_account a
left join sys_code c on a.status = c.sys_code and c.sys_code_type = 'user_job_status'
left join sys_code sex on a.sex = sex.sys_code and sex.sys_code_type = 'sex'
where ifnull(a.status, '1') != '0' and a.is_delete = '0'
require('./base_where.common')
