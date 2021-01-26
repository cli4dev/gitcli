package ui

const srcPublicEnvConfJson = `
{
    "name": "xxx系统",
    "copyright": {
        "company": "四川千行你我科技股份有限公司",
        "code": "蜀ICP备20003360号"
    },
    "system": {},
    "api": {
        "host": "http://localhost:8089",
        "verifyURL": "/sso/login/verify",
        "confURL": "/system/webconfig",
        "enumURL": "/dictionary/query",
        "logoutURL": "/sso/logout"
    },
    "sso": {
        "ident": "sso",
        "host": "http://ssov4.100bm0.com:6687"
    },
    "menus": [
        {
            "name": "日常管理",
            "icon": "-",
            "path": "-",
            "children": [
                {
                    "name": "交易管理",
                    "is_open": "1",
                    "icon": "fa fa-line-chart text-danger",
                    "path": "-",
                    "children": [
                        {
                            "name": "交易订单",
                            "icon": "fa fa-user-circle text-primary",
                            "path": "/order"
                        }
                    ]
                }
            ]
        }
    ],
    "sysList": []
}
`
