# Neoenergia Private API

- gerado em: `2026-04-16T22:19:31.650366`
- cliente: `Paula Pereira Fernandes`
- documento: `*******7586`
- artefatos: `docs/neoenergia-private-api/live_20260416_221852`

## Fluxo observado

1. Login no frontend da Agência Virtual.
2. Captura do Bearer token no `localStorage` após autenticação.
3. Consumo dos endpoints privados com `Authorization: Bearer ...`.
4. Obtenção de protocolo e chamadas por UC para conta, faturas, histórico e PDF.

## Seleção de estado

Nenhum endpoint dedicado de seleção de estado foi observado no fluxo capturado. Após o login, as chamadas já seguem com `distribuidora=COELBA`, o que indica que a escolha de estado é resolvida no frontend ou por configuração local de sessão.

## Unidades consumidoras observadas

- `007098175908` | status `LIGADA` | endereço `RUA BAHIA, 130` | município `LAPAO`
- `007085489032` | status `DESLIGADA` | endereço `RUA BAHIA, 130` | município `LAPAO`

## Endpoints documentados

### `GET /`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<!DOCTYPE html>\n<html lang=\"pt-br\">\n    <head>\n        <meta charset=\"utf-8\"/>\n        <title>Neoenergia - Agência Virtual</title>\n        <base href=\"/\"/>\n\n        <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0\"/>\n        <meta http-equiv=\"x-ua-compatible\" content=\"IE=edge\"/>\n        <link rel=\"icon\" type=\"image/x-icon\" href=\"assets/images/favicon.ico\"/>\n        <link rel=\"stylesheet\" type=\"text/css\" href=\"https://use.fontawesome.com/releases/v5.8.2/css/all.css\"/>\n\n        <!-- Google Tag Manager -->\n        <script>(function(w,d,s,l,i){w[l]=w[l]||[];w[l].push({'gtm.start':\n            new Date().getTime(),event:'gtm.js'});var f=d.getElementsByTagName(s)[0],\n            j=d.createElement(s),dl=l!='dataLayer'?'&l='+l:'';j.async=true;j.src=\n            'https://www.googletagmanager.com/gtm.js?id='+i+dl;f.parentNode.insertBefore(j,f);\n            })(window,document,'script','dataLayer','GTM-WQJP4MTG');\n        </script>\n        <!-- End Google Tag Manager -->\n\n        <script>\n            window.dataLayer = window.dataLayer || [];\n            function gtag() { dataLayer.push(arguments); }\n            gtag('js', new Date());\n        </script>\n\n    <link rel=\"stylesheet\" href=\"styles.b8f530f2278fd96d.css\"><script >bazadebezolkohpepadr=\"843664393\"</script><script type=\"text/javascript\" src=\"https://agenciavirtual.neoenergia.com/akam/13/32494caf\"  defer></script></head>\n\n    <body>\n        <!-- Google Tag Manager (noscript) -->\n        <noscript><iframe src=\"https://www.googletagmanager.com/ns.html?id=GTM-WQJP4MTG\"\n        height=\"0\" width=\"0\" style=\"display:none;visibility:hidden\"></iframe></noscript>\n        <!-- End Google Tag Manager (noscript) -->\n\n        <app-root></app-root>\n    <script src=\"runtime.a71dfaecea5c5d9f.js\" type=\"module\"></script><script src=\"polyfills.b50b1d0ce818639d.js\" type=\"module\"></script><script src=\"main.12d1bbaf63a33bd2.js\" type=\"module\"></script><script type=\"text/javascript\"  src=\"/Z21f/aWAX/7m/5GU3/k93g/3cw3mGrL4mGfJcG71L/Ay4KP0QVAw/IEB5E/QJ1ch4B\"></script><noscript><img src=\"https://agenciavirtual.neoenergia.com/akam/13/pixel_32494caf?a=dD0xY2EyNzMyMmE0NTM1ZmE1NWM1NTljMjcxNTZhNTQ5NmIxYmZiZTczJmpzPW9mZg==\" style=\"visibility: hidden; position: absolute; left: -999px; top: -999px;\" /></noscript></body>\n\n</html>\n"
```

### `GET /2002.6f882c2771439426.js`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"\"use strict\";\n(self[\"webpackChunkneoenergia_web\"] = self[\"webpackChunkneoenergia_web\"] || []).push([[2002],{\n\n/***/ 94222:\n/*!**********************************************************************************************************!*\\\n  !*** ./src/app/core/services/atendimento-presencial-agendado/atendimento-presencial-agendado.service.ts ***!\n  \\**********************************************************************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"AtendimentoPresencialAgendadoService\": () => (/* binding */ AtendimentoPresencialAgendadoService)\n/* harmony export */ });\n/* harmony import */ var _environments_environment__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @environments/environment */ 24766);\n/* harmony import */ var app_core_enums_distribuidoras__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! app/core/enums/distribuidoras */ 65700);\n/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @angular/core */ 5000);\n\n\n\nlet AtendimentoPresencialAgendadoService = /*#__PURE__*/(() => {\n  class AtendimentoPresencialAgendadoService {\n    getLinkAtendimentoPresencialAgendado() {\n      const hyperlinkMap = {\n        [app_core_enums_distribuidoras__WEBPACK_IMPORTED_MODULE_1__.Canal.AGE]: 'https://neoenergia.agendamento.ai/',\n        [app_core_enums_distribuidoras__WEBPACK_IMPORTED_MODULE_1__.Canal.AGR]: 'https://neoenergia.agendamento.ai/',\n        [app_core_enums_distribuidoras__WEBPACK_IMPORTED_MODULE_1__.Canal.AGC]: 'https://neoenergia.agendamento.ai/',\n        [app_core_enums_distribuidoras__WEBPACK_IMPORTED_MODULE_1__.Canal.AGP]: 'https://neoenergia.agendamento.ai/',\n        [app_core_enums_distribuidoras__WEBPACK_IMPORTED_MODULE_1__.Canal.AGU]: 'https://neoenergia.agendamento.ai/'\n      };\n      return hyperlinkMap[_environments_environment__WEBPACK_IMPORTED_MODULE_0__.environment.canal];\n    }\n\n  }\n\n  AtendimentoPresencialAgendadoService.ɵfac = function AtendimentoPresencialAgendadoService_Factory(t) {\n    return new (t || AtendimentoPresencialAgendadoService)();\n  };\n\n  AtendimentoPresencialAgendadoService.ɵprov = /*@__PURE__*/_angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵdefineInjectable\"]({\n    token: AtendimentoPresencialAgendadoService,\n    factory: AtendimentoPresencialAgendadoService.ɵfac,\n    providedIn: 'root'\n  });\n  return AtendimentoPresencialAgendadoService;\n})();\n\n/***/ }),\n\n/***/ 42002:\n/*!**************************************************************************!*\\\n  !*** ./src/app/modules/selecionar-estado/selecionar-estado.component.ts ***!\n  \\**************************************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"SelecionarEstadoComponent\": () => (/* binding */ SelecionarEstadoComponent)\n/* harmony export */ });\n/* harmony import */ var _environments_environment__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @environments/environment */ 24766);\n/* harmony import */ var app_core_enums_distribuidoras__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! app/core/enums/distribuidoras */ 65700);\n/* harmony import */ var app_core_enums_servicos__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! app/core/enums/servicos */ 68636);\n/* harmony import */ var app_core_interfaces_selecionar_estado__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! app/core/interfaces/selecionar-estado */ 69580);\n/* harmony import */ var app_core_models_RecuperarSenhaDTO_recuperarSenha__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! app/core/models/RecuperarSenhaDTO/recuperarSenha */ 11727);\n/* harmony import */ var _selecionar"
```

### `GET /2710.e1570061efda8ac2.js`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"\"use strict\";\n(self[\"webpackChunkneoenergia_web\"] = self[\"webpackChunkneoenergia_web\"] || []).push([[2710],{\n\n/***/ 52150:\n/*!**********************************************************************!*\\\n  !*** ./src/app/shared/components/pagination/pagination.component.ts ***!\n  \\**********************************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"PaginationComponent\": () => (/* binding */ PaginationComponent)\n/* harmony export */ });\n/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @angular/core */ 5000);\n/* harmony import */ var app_core_services_utils_neo_utils_service__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! app/core/services/utils/neo-utils.service */ 93225);\n/* harmony import */ var app_core_services_user_user_service__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! app/core/services/user/user.service */ 27353);\n/* harmony import */ var _angular_common__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! @angular/common */ 69808);\n/* harmony import */ var _ng_bootstrap_ng_bootstrap__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! @ng-bootstrap/ng-bootstrap */ 13707);\n\n\n\n\n\n\n\nfunction PaginationComponent_div_1_ng_template_2_Template(rf, ctx) {\n  if (rf & 1) {\n    _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵtext\"](0, \"<\");\n  }\n}\n\nfunction PaginationComponent_div_1_ng_template_3_Template(rf, ctx) {\n  if (rf & 1) {\n    _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵtext\"](0, \">\");\n  }\n}\n\nfunction PaginationComponent_div_1_Template(rf, ctx) {\n  if (rf & 1) {\n    const _r5 = _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵgetCurrentView\"]();\n\n    _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵelementStart\"](0, \"div\", 2)(1, \"ngb-pagination\", 3);\n    _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵlistener\"](\"pageChange\", function PaginationComponent_div_1_Template_ngb_pagination_pageChange_1_listener($event) {\n      _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵrestoreView\"](_r5);\n      const ctx_r4 = _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵnextContext\"]();\n      return ctx_r4.pageIndex = $event;\n    })(\"pageChange\", function PaginationComponent_div_1_Template_ngb_pagination_pageChange_1_listener($event) {\n      _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵrestoreView\"](_r5);\n      const ctx_r6 = _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵnextContext\"]();\n      return ctx_r6.onPageChanged($event);\n    });\n    _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵtemplate\"](2, PaginationComponent_div_1_ng_template_2_Template, 1, 0, \"ng-template\", 4);\n    _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵtemplate\"](3, PaginationComponent_div_1_ng_template_3_Template, 1, 0, \"ng-template\", 5);\n    _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵelementEnd\"]()();\n  }\n\n  if (rf & 2) {\n    const ctx_r0 = _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵnextContext\"]();\n    _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵclassProp\"](\"isGrupoA\", ctx_r0.user.group === \"A\")(\"isGrupoB\", ctx_r0.user.group === \"B\");\n    _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵadvance\"](1);\n    _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵproperty\"](\"collectionSize\", ctx_r0.collections.length)(\"page\", ctx_r0.pageIndex)(\"pageSize\", ctx_r0.pageSize)(\"ellipses\", false)(\"rotate\", true)(\"boundaryLinks\", true)(\"maxSize\", 2)(\"size\", \"sm\");\n  }\n}\n\nfunction PaginationComponent_div_2_ng_template_2_Template(rf, ctx) {\n  if (rf & 1) {\n    _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵtext\"](0, \"Anterior\");\n  }\n}\n\nfunction PaginationComponent_div_2_ng_template_3_Template(rf, ctx) {\n  if (rf & 1) {\n    _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵtext\"](0, \"Pr\\u00F3ximo\");\n  }\n}\n\nfunction PaginationComponent_div_2_4_ng_template_0_Template(rf, ctx) {\n  if (rf & 1) {\n    _angular_core__WEBPACK_IMPO"
```

### `GET /3597.7eb7cec1e3cc5e4a.js`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"\"use strict\";\n(self[\"webpackChunkneoenergia_web\"] = self[\"webpackChunkneoenergia_web\"] || []).push([[3597],{\n\n/***/ 39386:\n/*!****************************************************************************************!*\\\n  !*** ./src/app/core/routes/login-route/login-child-routes/selecao-de-perfis.routes.ts ***!\n  \\****************************************************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"selecaoPerfisRoutes\": () => (/* binding */ selecaoPerfisRoutes)\n/* harmony export */ });\n/* harmony import */ var app_modules_selecao_de_perfis_selecao_perfis_component__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! app/modules/selecao-de-perfis/selecao-perfis.component */ 16043);\n\nconst selecaoPerfisRoutes = [{\n  path: '',\n  component: app_modules_selecao_de_perfis_selecao_perfis_component__WEBPACK_IMPORTED_MODULE_0__.SelecaoPerfisComponent\n}];\n\n/***/ }),\n\n/***/ 7286:\n/*!*****************************************************************************************************************!*\\\n  !*** ./src/app/modules/selecao-de-perfis/components/dados-perfil-de-acesso/dados-perfil-de-acesso.component.ts ***!\n  \\*****************************************************************************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"DadosPerfilDeAcessoComponent\": () => (/* binding */ DadosPerfilDeAcessoComponent)\n/* harmony export */ });\n/* harmony import */ var app_core_enums_servicos__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! app/core/enums/servicos */ 68636);\n/* harmony import */ var app_core_models_home_sub_rotas_home__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! app/core/models/home/sub-rotas-home */ 32561);\n/* harmony import */ var app_core_models_multilogin_multilogin_acesso__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! app/core/models/multilogin/multilogin-acesso */ 22836);\n/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(/*! @angular/core */ 5000);\n/* harmony import */ var app_core_services_selecao_perfil_de_acesso_selecao_perfil_de_acesso_service__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! app/core/services/selecao-perfil-de-acesso/selecao-perfil-de-acesso.service */ 35436);\n/* harmony import */ var app_core_services_multilogin_acesso_multilogin_acesso_service__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! app/core/services/multilogin-acesso/multilogin-acesso.service */ 92940);\n/* harmony import */ var app_core_services_selecao_de_imovel_selecao_de_imovel_service__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! app/core/services/selecao-de-imovel/selecao-de-imovel.service */ 16855);\n/* harmony import */ var app_core_services_user_user_service__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(/*! app/core/services/user/user.service */ 27353);\n/* harmony import */ var _angular_router__WEBPACK_IMPORTED_MODULE_8__ = __webpack_require__(/*! @angular/router */ 74202);\n/* harmony import */ var _angular_common__WEBPACK_IMPORTED_MODULE_9__ = __webpack_require__(/*! @angular/common */ 69808);\n\n\n\n\n\n\n\n\n\n\n\nconst _c0 = function (a0) {\n  return {\n    \"color\": a0\n  };\n};\n\nfunction DadosPerfilDeAcessoComponent_div_1_Template(rf, ctx) {\n  if (rf & 1) {\n    const _r4 = _angular_core__WEBPACK_IMPORTED_MODULE_7__[\"ɵɵgetCurrentView\"]();\n\n    _angular_core__WEBPACK_IMPORTED_MODULE_7__[\"ɵɵelementStart\"](0, \"div\", 2)(1, \"button\", 3);\n    _angular_core__WEBPACK_IMPORTED_MODULE_7__[\"ɵɵlistener\"](\"click\", function DadosPerfilDeAcessoComponent_div_1_Template_button_click_1_listener() {\n      const restoredCtx = _angular_core__WEBPACK_IMPOR"
```

### `GET /3957.3daed7ce9f2584f5.js`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"\"use strict\";\n(self[\"webpackChunkneoenergia_web\"] = self[\"webpackChunkneoenergia_web\"] || []).push([[3957],{\n\n/***/ 37662:\n/*!*********************************************!*\\\n  !*** ./src/app/core/enums/dados-consumo.ts ***!\n  \\*********************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"EnumTipoConta\": () => (/* binding */ EnumTipoConta)\n/* harmony export */ });\nvar EnumTipoConta = /*#__PURE__*/(() => {\n  (function (EnumTipoConta) {\n    EnumTipoConta[\"CA\"] = \"Cativo\";\n    EnumTipoConta[\"CL\"] = \"Livre\";\n    EnumTipoConta[\"GR\"] = \"Gerador\";\n  })(EnumTipoConta || (EnumTipoConta = {}));\n\n  return EnumTipoConta;\n})();\n\n/***/ }),\n\n/***/ 63508:\n/*!**********************************************************************!*\\\n  !*** ./src/app/core/guards/meus-imoveis-guard/meus-imoveis.guard.ts ***!\n  \\**********************************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"MeusImoveisGuard\": () => (/* binding */ MeusImoveisGuard)\n/* harmony export */ });\n/* harmony import */ var app_core_enums_servicos__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! app/core/enums/servicos */ 68636);\n/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @angular/core */ 5000);\n/* harmony import */ var _angular_router__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! @angular/router */ 74202);\n/* harmony import */ var app_core_services_selecao_de_imovel_selecao_de_imovel_service__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! app/core/services/selecao-de-imovel/selecao-de-imovel.service */ 16855);\n\n\n\n\nlet MeusImoveisGuard = /*#__PURE__*/(() => {\n  class MeusImoveisGuard {\n    constructor(_router, _selecaoImovelService) {\n      this._router = _router;\n      this._selecaoImovelService = _selecaoImovelService;\n    }\n\n    canActivate() {\n      if (this._selecaoImovelService.getInformacoesUCSelecionada) {\n        return true;\n      }\n\n      this._router.navigate([app_core_enums_servicos__WEBPACK_IMPORTED_MODULE_0__.PathCompleto.meusImoveis]);\n\n      return false;\n    }\n\n  }\n\n  MeusImoveisGuard.ɵfac = function MeusImoveisGuard_Factory(t) {\n    return new (t || MeusImoveisGuard)(_angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵinject\"](_angular_router__WEBPACK_IMPORTED_MODULE_3__.Router), _angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵinject\"](app_core_services_selecao_de_imovel_selecao_de_imovel_service__WEBPACK_IMPORTED_MODULE_1__.SelecaoImovelService));\n  };\n\n  MeusImoveisGuard.ɵprov = /*@__PURE__*/_angular_core__WEBPACK_IMPORTED_MODULE_2__[\"ɵɵdefineInjectable\"]({\n    token: MeusImoveisGuard,\n    factory: MeusImoveisGuard.ɵfac,\n    providedIn: 'root'\n  });\n  return MeusImoveisGuard;\n})();\n\n/***/ }),\n\n/***/ 4958:\n/*!************************************************************************************************!*\\\n  !*** ./src/app/core/guards/minhas-unidades-consumidoras/minhas-unidades-consumidoras.guard.ts ***!\n  \\************************************************************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"MinhasUnidadesConsumidorasGuard\": () => (/* binding */ MinhasUnidadesConsumidorasGuard)\n/* harmony export */ });\n/* harmony import */ var app_core_enums_servicos__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! app/core/enums/servicos */ 68636);\n/* harmony import */ var app_core_models_multilogin_multilogin_acesso__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! app/co"
```

### `GET /4092.ad2a52232bdf2265.js`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"\"use strict\";\n(self[\"webpackChunkneoenergia_web\"] = self[\"webpackChunkneoenergia_web\"] || []).push([[4092],{\n\n/***/ 23343:\n/*!*********************************************************!*\\\n  !*** ./src/app/shared/Validators/password-validator.ts ***!\n  \\*********************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"PasswordValidator\": () => (/* binding */ PasswordValidator)\n/* harmony export */ });\nclass PasswordValidator {\n  static hasSequencialNumbers() {\n    return control => {\n      const regex = /^012|123|234|345|456|567|678|789|891|987|876|765|654|543|432|321|210$/;\n      const hasSequencialNumbers = regex.test(control.value);\n      return hasSequencialNumbers ? {\n        hasSequencialNumbers: true\n      } : null;\n    };\n  }\n\n  static hasPartsOfCPF(documento) {\n    return control => {\n      const numerosDocumento = documento.replace(/\\D/g, '');\n      const passwordArray = !!control.value ? control.value.split(/\\D/g) : '';\n      let hasPartsOfCPF = false;\n      let value = passwordArray[0] == null ? '' : passwordArray[0];\n\n      for (let i = 0; i <= numerosDocumento.length - 3; i++) {\n        const trecho = numerosDocumento.substring(i, i + 3);\n\n        if (value.includes(trecho)) {\n          hasPartsOfCPF = true;\n        }\n      }\n\n      return hasPartsOfCPF ? {\n        hasPartsOfCPF: true\n      } : null;\n    };\n  }\n\n  static hasRepeatParts() {\n    return control => {\n      const regex = /([0-9]+)\\1{1,}/;\n      const hasRepeatParts = regex.test(control.value);\n      return hasRepeatParts ? {\n        hasRepeatParts: true\n      } : null;\n    };\n  }\n\n  static hasNotAllowedCharacters() {\n    return control => {\n      const regex = /^[a-zA-Z0-9!?@#~^ ]+$/;\n      const hasNotAllowedCharacters = !!control.value && !regex.test(control.value);\n      return hasNotAllowedCharacters ? {\n        hasNotAllowedCharacters: true\n      } : null;\n    };\n  }\n\n  static hasEmptyString() {\n    return control => {\n      let hasEmptyString = false;\n      if (!!control.value) hasEmptyString = control.value.includes(\" \");\n      return hasEmptyString ? {\n        hasEmptyString: true\n      } : null;\n    };\n  }\n\n}\n\n/***/ }),\n\n/***/ 90592:\n/*!******************************************************************************!*\\\n  !*** ./src/app/shared/components/error-criteria/error-criteria.component.ts ***!\n  \\******************************************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"ErrorCriteriaComponent\": () => (/* binding */ ErrorCriteriaComponent)\n/* harmony export */ });\n/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ 5000);\n/* harmony import */ var _angular_common__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/common */ 69808);\n/* harmony import */ var _angular_material_icon__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @angular/material/icon */ 25245);\n\n\n\n\nconst _c0 = function (a0) {\n  return {\n    \"fw-bold\": a0\n  };\n};\n\nlet ErrorCriteriaComponent = /*#__PURE__*/(() => {\n  class ErrorCriteriaComponent {\n    constructor() {\n      this.errorName = \"\";\n      this.message = \"\";\n    }\n\n  }\n\n  ErrorCriteriaComponent.ɵfac = function ErrorCriteriaComponent_Factory(t) {\n    return new (t || ErrorCriteriaComponent)();\n  };\n\n  ErrorCriteriaComponent.ɵcmp = /*@__PURE__*/_angular_core__WEBPACK_IMPORTED_MODULE_0__[\"ɵɵdefineComponent\"]({\n    type: ErrorCriteriaComponent,\n    selectors: [[\"app-error-criteria\"]],\n    inputs: {\n      passwordFormGroup: \"passwordFormGroup\",\n      errorName: \"errorName\",\n      message: \"message\"\n    },\n    decls: 6,\n    vars"
```

### `GET /4260.09709e47c7d90021.js`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"\"use strict\";\n(self[\"webpackChunkneoenergia_web\"] = self[\"webpackChunkneoenergia_web\"] || []).push([[4260],{\n\n/***/ 62500:\n/*!***************************************************************!*\\\n  !*** ./src/app/core/models/multilogin/multilogin-cadastro.ts ***!\n  \\***************************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"AcessoComumCompartilhaCom\": () => (/* binding */ AcessoComumCompartilhaCom),\n/* harmony export */   \"ConjugeCompartilhaCom\": () => (/* binding */ ConjugeCompartilhaCom),\n/* harmony export */   \"CredenciadoCompartilhaCom\": () => (/* binding */ CredenciadoCompartilhaCom),\n/* harmony export */   \"ImobiliariaCompartilhaCom\": () => (/* binding */ ImobiliariaCompartilhaCom),\n/* harmony export */   \"InformacoesUsuario\": () => (/* binding */ InformacoesUsuario),\n/* harmony export */   \"MultiloginCompartilharAcesso\": () => (/* binding */ MultiloginCompartilharAcesso),\n/* harmony export */   \"RepresetanteLegalCompartilhaCom\": () => (/* binding */ RepresetanteLegalCompartilhaCom),\n/* harmony export */   \"SubRotasMultiloginCadastro\": () => (/* binding */ SubRotasMultiloginCadastro),\n/* harmony export */   \"ValidaRelacao\": () => (/* binding */ ValidaRelacao),\n/* harmony export */   \"emailsBO\": () => (/* binding */ emailsBO),\n/* harmony export */   \"servicosMKTAutomation\": () => (/* binding */ servicosMKTAutomation)\n/* harmony export */ });\n/* harmony import */ var _documentos_box_anexo_box_anexo__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ../documentos/box-anexo/box-anexo */ 78014);\n/* harmony import */ var _multilogin_acesso__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./multilogin-acesso */ 22836);\n\n\nvar SubRotasMultiloginCadastro = /*#__PURE__*/(() => {\n  (function (SubRotasMultiloginCadastro) {\n    SubRotasMultiloginCadastro[\"Imobiliaria\"] = \"imobiliario\";\n    SubRotasMultiloginCadastro[\"Credenciado\"] = \"estabelecimento-credenciado\";\n    SubRotasMultiloginCadastro[\"Compartilhamento\"] = \"compartilhar-acesso\";\n    SubRotasMultiloginCadastro[\"Avisos\"] = \"avisos\";\n    SubRotasMultiloginCadastro[\"ImobiliariaANL\"] = \"cadastro-imobiliario\";\n    SubRotasMultiloginCadastro[\"CredenciadoANL\"] = \"cadastro-credenciado\";\n    SubRotasMultiloginCadastro[\"CadastroDeParceiros\"] = \"cadastro-de-parceiros\";\n  })(SubRotasMultiloginCadastro || (SubRotasMultiloginCadastro = {}));\n\n  return SubRotasMultiloginCadastro;\n})();\nclass MultiloginCompartilharAcesso {\n  constructor(tipoAtribuicao, nomeFiscalSecundario, documentoFiscalSecundario, dataVigenciaContrato, boxAnexo, comprovantes, tipoEmailBO, customerKey) {\n    this.tipoAtribuicao = tipoAtribuicao;\n    this.nomeFiscalSecundario = nomeFiscalSecundario;\n    this.documentoFiscalSecundario = documentoFiscalSecundario;\n    this.dataVigenciaContrato = dataVigenciaContrato;\n    this.boxAnexo = boxAnexo;\n    this.comprovantes = comprovantes;\n    this.tipoEmailBO = tipoEmailBO;\n    this.customerKey = customerKey;\n    this.boxAnexo = new _documentos_box_anexo_box_anexo__WEBPACK_IMPORTED_MODULE_0__.BoxAnexo('DOCUMENTO DE COMPROVAÇÃO', false, 'DOCUMENTO DE COMPROVAÇÃO');\n    this.comprovantes = [];\n  }\n\n}\nclass InformacoesUsuario {\n  constructor(cpf, nome, inicio, fim, button, isValid) {\n    this.cpf = cpf;\n    this.nome = nome;\n    this.inicio = inicio;\n    this.fim = fim;\n    this.button = button;\n    this.isValid = isValid;\n  }\n\n}\nconst ValidaRelacao = [{\n  key: 'C',\n  value: _multilogin_acesso__WEBPACK_IMPORTED_MODULE_1__.PerfisDeAcesso.conjuge\n}, {\n  key: 'R',\n  value: _multilogin_acesso__WEBPACK_IMPORTED_MODULE_1__.PerfisDeAcesso.representanteLegal\n}];\nconst AcessoComumCompartilhaCom = [{\n  key: 'ACESSO COMUM COM REPRESENTANTE LEGAL',\n  value: _multilogin_acesso__WEBPACK_IMPORTED_MODULE_1__.PerfisDeAcesso.representanteLegal\n}, {\n  key: 'ACESSO COMUM COM PADRO"
```

### `GET /6856.c1b940162fba60a2.js`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"\"use strict\";\n(self[\"webpackChunkneoenergia_web\"] = self[\"webpackChunkneoenergia_web\"] || []).push([[6856],{\n\n/***/ 86856:\n/*!****************************************************************!*\\\n  !*** ./node_modules/@angular/material/fesm2015/datepicker.mjs ***!\n  \\****************************************************************/\n/***/ ((__unused_webpack___webpack_module__, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"DateRange\": () => (/* binding */ DateRange),\n/* harmony export */   \"DefaultMatCalendarRangeStrategy\": () => (/* binding */ DefaultMatCalendarRangeStrategy),\n/* harmony export */   \"MAT_DATEPICKER_SCROLL_STRATEGY\": () => (/* binding */ MAT_DATEPICKER_SCROLL_STRATEGY),\n/* harmony export */   \"MAT_DATEPICKER_SCROLL_STRATEGY_FACTORY\": () => (/* binding */ MAT_DATEPICKER_SCROLL_STRATEGY_FACTORY),\n/* harmony export */   \"MAT_DATEPICKER_SCROLL_STRATEGY_FACTORY_PROVIDER\": () => (/* binding */ MAT_DATEPICKER_SCROLL_STRATEGY_FACTORY_PROVIDER),\n/* harmony export */   \"MAT_DATEPICKER_VALIDATORS\": () => (/* binding */ MAT_DATEPICKER_VALIDATORS),\n/* harmony export */   \"MAT_DATEPICKER_VALUE_ACCESSOR\": () => (/* binding */ MAT_DATEPICKER_VALUE_ACCESSOR),\n/* harmony export */   \"MAT_DATE_RANGE_SELECTION_STRATEGY\": () => (/* binding */ MAT_DATE_RANGE_SELECTION_STRATEGY),\n/* harmony export */   \"MAT_RANGE_DATE_SELECTION_MODEL_FACTORY\": () => (/* binding */ MAT_RANGE_DATE_SELECTION_MODEL_FACTORY),\n/* harmony export */   \"MAT_RANGE_DATE_SELECTION_MODEL_PROVIDER\": () => (/* binding */ MAT_RANGE_DATE_SELECTION_MODEL_PROVIDER),\n/* harmony export */   \"MAT_SINGLE_DATE_SELECTION_MODEL_FACTORY\": () => (/* binding */ MAT_SINGLE_DATE_SELECTION_MODEL_FACTORY),\n/* harmony export */   \"MAT_SINGLE_DATE_SELECTION_MODEL_PROVIDER\": () => (/* binding */ MAT_SINGLE_DATE_SELECTION_MODEL_PROVIDER),\n/* harmony export */   \"MatCalendar\": () => (/* binding */ MatCalendar),\n/* harmony export */   \"MatCalendarBody\": () => (/* binding */ MatCalendarBody),\n/* harmony export */   \"MatCalendarCell\": () => (/* binding */ MatCalendarCell),\n/* harmony export */   \"MatCalendarHeader\": () => (/* binding */ MatCalendarHeader),\n/* harmony export */   \"MatDateRangeInput\": () => (/* binding */ MatDateRangeInput),\n/* harmony export */   \"MatDateRangePicker\": () => (/* binding */ MatDateRangePicker),\n/* harmony export */   \"MatDateSelectionModel\": () => (/* binding */ MatDateSelectionModel),\n/* harmony export */   \"MatDatepicker\": () => (/* binding */ MatDatepicker),\n/* harmony export */   \"MatDatepickerActions\": () => (/* binding */ MatDatepickerActions),\n/* harmony export */   \"MatDatepickerApply\": () => (/* binding */ MatDatepickerApply),\n/* harmony export */   \"MatDatepickerCancel\": () => (/* binding */ MatDatepickerCancel),\n/* harmony export */   \"MatDatepickerContent\": () => (/* binding */ MatDatepickerContent),\n/* harmony export */   \"MatDatepickerInput\": () => (/* binding */ MatDatepickerInput),\n/* harmony export */   \"MatDatepickerInputEvent\": () => (/* binding */ MatDatepickerInputEvent),\n/* harmony export */   \"MatDatepickerIntl\": () => (/* binding */ MatDatepickerIntl),\n/* harmony export */   \"MatDatepickerModule\": () => (/* binding */ MatDatepickerModule),\n/* harmony export */   \"MatDatepickerToggle\": () => (/* binding */ MatDatepickerToggle),\n/* harmony export */   \"MatDatepickerToggleIcon\": () => (/* binding */ MatDatepickerToggleIcon),\n/* harmony export */   \"MatEndDate\": () => (/* binding */ MatEndDate),\n/* harmony export */   \"MatMonthView\": () => (/* binding */ MatMonthView),\n/* harmony export */   \"MatMultiYearView\": () => (/* binding */ MatMultiYearView),\n/* harmony export */   \"MatRangeDateSelectionModel\": () => (/* binding */ MatRangeDateSelectionModel),\n/* harmony export */   \"MatSingleDateSelectionModel\": () => (/* binding */ MatSingleDateSelectionModel),\n/* harmony export */   \"MatStartDate\": () => (/* binding"
```

### `GET /7618.38f40eb2d1cd74df.js`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"\"use strict\";\n(self[\"webpackChunkneoenergia_web\"] = self[\"webpackChunkneoenergia_web\"] || []).push([[7618],{\n\n/***/ 45120:\n/*!******************************************************!*\\\n  !*** ./src/app/core/models/flex-pag/flex-pag-dto.ts ***!\n  \\******************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"ClienteDTO\": () => (/* binding */ ClienteDTO),\n/* harmony export */   \"Debts\": () => (/* binding */ Debts),\n/* harmony export */   \"Endereco\": () => (/* binding */ Endereco),\n/* harmony export */   \"FlexPagRequest\": () => (/* binding */ FlexPagRequest),\n/* harmony export */   \"FlexPagResponse\": () => (/* binding */ FlexPagResponse),\n/* harmony export */   \"FlexPagTipificacao\": () => (/* binding */ FlexPagTipificacao),\n/* harmony export */   \"InvoiceDTO\": () => (/* binding */ InvoiceDTO),\n/* harmony export */   \"Order\": () => (/* binding */ Order),\n/* harmony export */   \"Principal\": () => (/* binding */ Principal),\n/* harmony export */   \"UC\": () => (/* binding */ UC)\n/* harmony export */ });\n/* harmony import */ var app_core_enums_regiao__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! app/core/enums/regiao */ 35343);\n\nvar FlexPagTipificacao = /*#__PURE__*/(() => {\n  (function (FlexPagTipificacao) {\n    FlexPagTipificacao[\"tipificacao\"] = \"1010810\";\n  })(FlexPagTipificacao || (FlexPagTipificacao = {}));\n\n  return FlexPagTipificacao;\n})();\nclass FlexPagResponse {\n  constructor() {\n    this.url = \"\";\n  }\n\n}\nclass FlexPagRequest {\n  constructor() {\n    this.tipoPerfil = 1;\n    this.is_authenticated = \"true\"; // valor fixo\n\n    this.is_principal_hide = \"true\"; // valor fixo\n\n    this.payment_module = \"INVOICES_CREATE\";\n    this.origin = 'NEOENERGIA_AGENCIA';\n    this.canalSolicitante = 'AGC';\n    this.usuario = 'UCSCOMM';\n    this.regiao = app_core_enums_regiao__WEBPACK_IMPORTED_MODULE_0__.Regiao.NE;\n    this.distribuidora = 'COELBA';\n    this.protocolo = '';\n    this.tipificacao = FlexPagTipificacao.tipificacao;\n    this.numeroFatura = '';\n    this.recaptchaAnl = false;\n    this.recaptcha = '';\n    this.order = new Order();\n    this.codigo = '';\n    this.dataFimVencFat = \"\";\n  }\n\n}\nclass Order {\n  constructor() {\n    this.id = '';\n    this.principal = new Principal();\n    this.debts = [];\n  }\n\n}\nclass Principal {\n  constructor() {\n    this.debt_source = 'COELBA';\n    this.document = '';\n    this.contract_account = '';\n  }\n\n}\nclass Debts {\n  constructor() {\n    this.invoice_id = '';\n    this.amount = '';\n    this.due_date = '';\n    this.barcode = '';\n  }\n\n}\nclass ClienteDTO {\n  constructor() {\n    this.nome = \"\";\n    this.tipoDocumento = \"\";\n    this.ni = \"\";\n    this.email = \"\";\n    this.telefone = \"\";\n    this.dataNascimento = \"\";\n    this.endereco = new Endereco();\n  }\n\n}\nclass Endereco {\n  constructor() {\n    this.bairro = \"\";\n    this.cep = \"\";\n    this.complemento = \"\";\n    this.logradouro = \"\";\n    this.municipio = \"\";\n    this.numero = \"\";\n    this.uf = \"\";\n  }\n\n}\nclass UC {\n  constructor() {\n    this.uc = \"\";\n    this.invoices = [];\n  }\n\n}\nclass InvoiceDTO {\n  constructor() {\n    this.amount = 0;\n    this.dueDate = \"\";\n    this.invoiceId = \"\";\n    this.barCodeOne = \"\";\n    this.barCodeTwo = \"\";\n  }\n\n}\n\n/***/ }),\n\n/***/ 29593:\n/*!******************************************************************!*\\\n  !*** ./src/app/core/models/multilogin/request/multilogin-dto.ts ***!\n  \\******************************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"AlterarOptinNotificacoesRequest\": () => (/* binding */ AlterarOptinNotificacoesRequest),\n/* harmony export */   \"AnexosDTORequest\": () => (/* binding */ AnexosDTOR"
```

### `GET /8026.b32fff5437314150.js`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"\"use strict\";\n(self[\"webpackChunkneoenergia_web\"] = self[\"webpackChunkneoenergia_web\"] || []).push([[8026],{\n\n/***/ 90806:\n/*!*************************************************************************************!*\\\n  !*** ./src/app/core/routes/login-route/login-child-routes/selecao-estado.routes.ts ***!\n  \\*************************************************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"selecaoEstadoRoutes\": () => (/* binding */ selecaoEstadoRoutes)\n/* harmony export */ });\n/* harmony import */ var app_modules_selecionar_estado_selecionar_estado_component__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! app/modules/selecionar-estado/selecionar-estado.component */ 42002);\n\nconst selecaoEstadoRoutes = [{\n  path: '',\n  component: app_modules_selecionar_estado_selecionar_estado_component__WEBPACK_IMPORTED_MODULE_0__.SelecionarEstadoComponent\n}];\n\n/***/ }),\n\n/***/ 58026:\n/*!***********************************************************************!*\\\n  !*** ./src/app/modules/selecionar-estado/selecionar-estado.module.ts ***!\n  \\***********************************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"SelecionarEstadoModule\": () => (/* binding */ SelecionarEstadoModule)\n/* harmony export */ });\n/* harmony import */ var app_shared_components_protocolo_informativo_protocolo_informativo_module__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! app/shared/components/protocolo-informativo/protocolo-informativo.module */ 42027);\n/* harmony import */ var _angular_common__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! @angular/common */ 69808);\n/* harmony import */ var _angular_material_card__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! @angular/material/card */ 9224);\n/* harmony import */ var _selecionar_estado_component__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./selecionar-estado.component */ 42002);\n/* harmony import */ var _angular_router__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(/*! @angular/router */ 74202);\n/* harmony import */ var app_core_routes_login_route_login_child_routes_selecao_estado_routes__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! app/core/routes/login-route/login-child-routes/selecao-estado.routes */ 90806);\n/* harmony import */ var _angular_material_icon__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(/*! @angular/material/icon */ 25245);\n/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! @angular/core */ 5000);\n\n\n\n\n\n\n\n\n\nlet SelecionarEstadoModule = /*#__PURE__*/(() => {\n  class SelecionarEstadoModule {}\n\n  SelecionarEstadoModule.ɵfac = function SelecionarEstadoModule_Factory(t) {\n    return new (t || SelecionarEstadoModule)();\n  };\n\n  SelecionarEstadoModule.ɵmod = /*@__PURE__*/_angular_core__WEBPACK_IMPORTED_MODULE_3__[\"ɵɵdefineNgModule\"]({\n    type: SelecionarEstadoModule\n  });\n  SelecionarEstadoModule.ɵinj = /*@__PURE__*/_angular_core__WEBPACK_IMPORTED_MODULE_3__[\"ɵɵdefineInjector\"]({\n    imports: [[_angular_common__WEBPACK_IMPORTED_MODULE_4__.CommonModule, _angular_material_card__WEBPACK_IMPORTED_MODULE_5__.MatCardModule, _angular_material_icon__WEBPACK_IMPORTED_MODULE_6__.MatIconModule, _angular_router__WEBPACK_IMPORTED_MODULE_7__.RouterModule.forChild(app_core_routes_login_route_login_child_routes_selecao_estado_routes__WEBPACK_IMPORTED_MODULE_2__.selecaoEstadoRoutes), app_shared_components_protocolo_informativo_protocolo_informativo_module__WEBPACK_IMPORTED_MODULE_0__.ProtocoloInformativoComponentModule]]\n  });\n  return SelecionarEstadoModule;\n})();\n\n(function () {\n  (typeof ngJitMode ="
```

### `GET /8970.44651cd0a12bb43c.js`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"\"use strict\";\n(self[\"webpackChunkneoenergia_web\"] = self[\"webpackChunkneoenergia_web\"] || []).push([[8970],{\n\n/***/ 59169:\n/*!**********************************************************************!*\\\n  !*** ./src/app/core/components/guia-dialog/guia-dialog.component.ts ***!\n  \\**********************************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"GuiaDialogComponent\": () => (/* binding */ GuiaDialogComponent)\n/* harmony export */ });\n/* harmony import */ var _angular_material_stepper__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! @angular/material/stepper */ 55615);\n/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/core */ 5000);\n/* harmony import */ var _angular_material_dialog__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @angular/material/dialog */ 48966);\n/* harmony import */ var _angular_material_button__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! @angular/material/button */ 47423);\n/* harmony import */ var _angular_material_icon__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! @angular/material/icon */ 25245);\n/* harmony import */ var _shared_components_neo_button_neo_button_component__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ../../../shared/components/neo-button/neo-button.component */ 56864);\n\n\n\n\n\n\n\nlet GuiaDialogComponent = /*#__PURE__*/(() => {\n  class GuiaDialogComponent {\n    constructor(dialogRef) {\n      this.dialogRef = dialogRef;\n      this.step = 1;\n      this.total = 4;\n    }\n\n    next() {\n      this.stepper.next();\n      this.step++;\n    }\n\n    previous() {\n      this.stepper.previous();\n      this.step--;\n    }\n\n    close(querCadastrar, terminouTutorial) {\n      this.dialogRef.close({\n        querCadastrar,\n        terminouTutorial\n      });\n    }\n\n  }\n\n  GuiaDialogComponent.ɵfac = function GuiaDialogComponent_Factory(t) {\n    return new (t || GuiaDialogComponent)(_angular_core__WEBPACK_IMPORTED_MODULE_1__[\"ɵɵdirectiveInject\"](_angular_material_dialog__WEBPACK_IMPORTED_MODULE_2__.MatDialogRef));\n  };\n\n  GuiaDialogComponent.ɵcmp = /*@__PURE__*/_angular_core__WEBPACK_IMPORTED_MODULE_1__[\"ɵɵdefineComponent\"]({\n    type: GuiaDialogComponent,\n    selectors: [[\"app-guia-dialog\"]],\n    viewQuery: function GuiaDialogComponent_Query(rf, ctx) {\n      if (rf & 1) {\n        _angular_core__WEBPACK_IMPORTED_MODULE_1__[\"ɵɵviewQuery\"](_angular_material_stepper__WEBPACK_IMPORTED_MODULE_3__.MatStepper, 5);\n      }\n\n      if (rf & 2) {\n        let _t;\n\n        _angular_core__WEBPACK_IMPORTED_MODULE_1__[\"ɵɵqueryRefresh\"](_t = _angular_core__WEBPACK_IMPORTED_MODULE_1__[\"ɵɵloadQuery\"]()) && (ctx.stepper = _t.first);\n      }\n    },\n    decls: 76,\n    vars: 23,\n    consts: [[3, \"linear\"], [\"stepper\", \"\"], [\"1\", \"\"], [1, \"guia-dialog-header\"], [1, \"guia-dialog-stepper\"], [\"type\", \"button\", \"mat-icon-button\", \"\", \"aria-label\", \"Fechar\", \"title\", \"Fechar\", 1, \"\", 3, \"click\"], [1, \"guia-dialog-content\"], [\"width\", \"250\", \"height\", \"250\", \"src\", \"/assets/images/pagina-inicial-1.png\", \"alt\", \"\"], [1, \"guia-dialog-footer\"], [3, \"classes\", \"titleText\", \"click\"], [\"2\", \"\"], [\"width\", \"200\", \"height\", \"200\", \"src\", \"/assets/images/pagina-inicial-2.png\", \"alt\", \"\"], [1, \"guia-dialog-link\"], [3, \"click\"], [\"3\", \"\"], [\"width\", \"250\", \"height\", \"250\", \"src\", \"/assets/images/pagina-inicial-3.png\", \"alt\", \"\"], [\"4\", \"\"], [\"width\", \"505\", \"height\", \"205\", \"src\", \"/assets/images/pagina-inicial-4.png\", \"alt\", \"\"]],\n    template: function GuiaDialogComponent_Template(rf, ctx) {\n      if (rf & 1) {\n        _angular_core__WEBPACK_IMPORTED_MODULE_1__[\"ɵɵelementStart\"](0, \"mat-stepper\", 0, 1)(2, \"mat-step\", null, 2)(4, \"div\", 3)(5, \"div\", 4)(6, \"p\");\n        _angular_core__WEBPACK_IMPORTED_MODULE_1__[\"ɵɵtext\"](7);\n        _angular_core__WEBPACK_"
```

### `GET /9831.cb12b7924c78a649.js`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"\"use strict\";\n(self[\"webpackChunkneoenergia_web\"] = self[\"webpackChunkneoenergia_web\"] || []).push([[9831],{\n\n/***/ 61556:\n/*!**********************************!*\\\n  !*** ./src/app/app.constants.ts ***!\n  \\**********************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"Constants\": () => (/* binding */ Constants)\n/* harmony export */ });\nlet Constants = /*#__PURE__*/(() => {\n  class Constants {}\n\n  Constants.userRole = {\n    Administrator: 'ROLE_ADMIN',\n    CommonUser: 'ROLE_COMMON_USER'\n  };\n  Constants.userRoles = [{\n    code: 'ROLE_ADMIN',\n    description: 'Administrator'\n  }];\n  return Constants;\n})();\n\n/***/ }),\n\n/***/ 19831:\n/*!******************************************************!*\\\n  !*** ./src/app/core/services/utils/utils.service.ts ***!\n  \\******************************************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"MONTHS\": () => (/* binding */ MONTHS),\n/* harmony export */   \"PATHS\": () => (/* binding */ PATHS),\n/* harmony export */   \"checkAndConvertExcelDateToJSDate\": () => (/* binding */ checkAndConvertExcelDateToJSDate),\n/* harmony export */   \"convertDates\": () => (/* binding */ convertDates),\n/* harmony export */   \"convertHours\": () => (/* binding */ convertHours),\n/* harmony export */   \"convertMonthDaysList\": () => (/* binding */ convertMonthDaysList),\n/* harmony export */   \"convertWeekList\": () => (/* binding */ convertWeekList),\n/* harmony export */   \"formatDateArrayToString\": () => (/* binding */ formatDateArrayToString),\n/* harmony export */   \"formatMonetaryNumber\": () => (/* binding */ formatMonetaryNumber),\n/* harmony export */   \"formatNumber\": () => (/* binding */ formatNumber),\n/* harmony export */   \"formatTooltip\": () => (/* binding */ formatTooltip),\n/* harmony export */   \"formatarMoeda\": () => (/* binding */ formatarMoeda),\n/* harmony export */   \"formatarStatus\": () => (/* binding */ formatarStatus),\n/* harmony export */   \"getDate\": () => (/* binding */ getDate),\n/* harmony export */   \"getEnderecoCompleto\": () => (/* binding */ getEnderecoCompleto),\n/* harmony export */   \"getFormattedCelular\": () => (/* binding */ getFormattedCelular),\n/* harmony export */   \"getFormattedTelefone\": () => (/* binding */ getFormattedTelefone),\n/* harmony export */   \"getKeysDayMonthFromJson\": () => (/* binding */ getKeysDayMonthFromJson),\n/* harmony export */   \"getKeysYearWeekFromJson\": () => (/* binding */ getKeysYearWeekFromJson),\n/* harmony export */   \"msgSituacao\": () => (/* binding */ msgSituacao),\n/* harmony export */   \"removeFirstAndLastCharacter\": () => (/* binding */ removeFirstAndLastCharacter),\n/* harmony export */   \"retornaDataFormatada\": () => (/* binding */ retornaDataFormatada),\n/* harmony export */   \"showActionAdministrator\": () => (/* binding */ showActionAdministrator),\n/* harmony export */   \"showDownloadsButtonsToSAP\": () => (/* binding */ showDownloadsButtonsToSAP),\n/* harmony export */   \"showMenuAdministrator\": () => (/* binding */ showMenuAdministrator),\n/* harmony export */   \"showMenuParametrization\": () => (/* binding */ showMenuParametrization),\n/* harmony export */   \"showMenuPerfomAlignment\": () => (/* binding */ showMenuPerfomAlignment),\n/* harmony export */   \"showMenuViewReport\": () => (/* binding */ showMenuViewReport),\n/* harmony export */   \"showMenuYourReference\": () => (/* binding */ showMenuYourReference),\n/* harmony export */   \"showPipelineOptimization\": () => (/* binding */ showPipelineOptimization),\n/* harmony export */   \"toDate\": () => (/* binding */ toDate),\n/* harmony export */   \"toDateTime\": () => (/* binding */ toDateTime),\n/* harmony export */   \"toLocalDate\": () => (/* bindi"
```

### `GET /MaterialIcons-Regular.7ea2023eeca07427.woff2`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<unreadable: 'utf-8' codec can't decode byte 0xad in position 10: invalid start byte>"
```

### `GET /Z21f/aWAX/7m/5GU3/k93g/3cw3mGrL4mGfJcG71L/Ay4KP0QVAw/IEB5E/QJ1ch4B`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"(function(){if(typeof Array.prototype.entries!=='function'){Object.defineProperty(Array.prototype,'entries',{value:function(){var index=0;const array=this;return {next:function(){if(index<array.length){return {value:[index,array[index++]],done:false};}else{return {done:true};}},[Symbol.iterator]:function(){return this;}};},writable:true,configurable:true});}}());(function(){X6();Acv();dTv();var Kn=function(){return wx.apply(this,[Sg,arguments]);};var Qj=function(XG){return +XG;};var O9=function(UX){return -UX;};var rj=function(nc,VP){return nc-VP;};var I3=function(zg,Kx){return zg|Kx;};function dTv(){Vh=fp+fp*Ov+XM*Ov*Ov+Ov*Ov*Ov,FU=XM+r1*Ov+gF*Ov*Ov+Ov*Ov*Ov,Hd=rd+MD*Ov+Ov*Ov,U3=rt+MD*Ov+Ov*Ov+Ov*Ov*Ov,t6=ZD+rt*Ov,Qh=r1+rt*Ov+Ov*Ov+Ov*Ov*Ov,d7=rd+XM*Ov+fp*Ov*Ov+Ov*Ov*Ov,NC=ZD+MD*Ov+Ov*Ov,Il=MD+rt*Ov+ZD*Ov*Ov,k7=rd+xf*Ov+fp*Ov*Ov+Ov*Ov*Ov,q=ZD+Ov+gF*Ov*Ov,pY=MD+gF*Ov+r1*Ov*Ov,NB=rt+rd*Ov+MD*Ov*Ov+Ov*Ov*Ov,kx=r1+XM*Ov+r1*Ov*Ov,B3=fp+xf*Ov+ZD*Ov*Ov+Ov*Ov*Ov,vh=gF+xf*Ov+fp*Ov*Ov+Ov*Ov*Ov,Hg=ZD+MD*Ov+Ov*Ov+Ov*Ov*Ov,LY=MD+rd*Ov,fT=N6+rd*Ov+fp*Ov*Ov+Ov*Ov*Ov,FT=ZD+Ov+MD*Ov*Ov+Ov*Ov*Ov,FS=ZD+gF*Ov+rd*Ov*Ov+Ov*Ov*Ov,qt=rd+MD*Ov+gF*Ov*Ov,rf=r1+xf*Ov+MD*Ov*Ov,kg=gF+rt*Ov+ZD*Ov*Ov+Ov*Ov*Ov,LB=r1+xf*Ov+XM*Ov*Ov+Ov*Ov*Ov,mS=xf+fp*Ov+fp*Ov*Ov+Ov*Ov*Ov,xv=r1+rd*Ov+Ov*Ov+Ov*Ov*Ov,GG=rt+MD*Ov+XM*Ov*Ov+Ov*Ov*Ov,Oj=N6+Ov+Ov*Ov+Ov*Ov*Ov,j2=XM+xf*Ov+MD*Ov*Ov+Ov*Ov*Ov,Or=ZD+Ov+fp*Ov*Ov,NU=MD+fp*Ov+rd*Ov*Ov,GC=rd+rt*Ov+ZD*Ov*Ov+r1*Ov*Ov*Ov+r1*Ov*Ov*Ov*Ov,KU=fp+rd*Ov+fp*Ov*Ov+Ov*Ov*Ov,O=rt+ZD*Ov+fp*Ov*Ov,lf=xf+rd*Ov+ZD*Ov*Ov,BU=fp+r1*Ov+r1*Ov*Ov+Ov*Ov*Ov,ZB=fp+Ov+XM*Ov*Ov+Ov*Ov*Ov,Ct=N6+r1*Ov+Ov*Ov,U8=XM+MD*Ov+ZD*Ov*Ov,B1=N6+xf*Ov+rd*Ov*Ov,kf=gF+Ov+Ov*Ov,rh=xf+Ov+gF*Ov*Ov+Ov*Ov*Ov,E9=gF+Ov+fp*Ov*Ov+Ov*Ov*Ov,Qx=rd+ZD*Ov+XM*Ov*Ov+Ov*Ov*Ov,Vr=XM+fp*Ov+gF*Ov*Ov,wP=MD+gF*Ov,Rd=ZD+gF*Ov+Ov*Ov,G3=MD+ZD*Ov+fp*Ov*Ov+Ov*Ov*Ov,Hx=XM+XM*Ov+Ov*Ov+Ov*Ov*Ov,sb=rt+rd*Ov+gF*Ov*Ov+Ov*Ov*Ov,lX=XM+gF*Ov+MD*Ov*Ov+Ov*Ov*Ov,Lc=ZD+gF*Ov+gF*Ov*Ov,V6=xf+XM*Ov+Ov*Ov,R9=fp+Ov+r1*Ov*Ov+Ov*Ov*Ov,bd=N6+xf*Ov+gF*Ov*Ov,bY=gF+fp*Ov,Bj=gF+Ov+rd*Ov*Ov+Ov*Ov*Ov,z7=r1+Ov+gF*Ov*Ov+Ov*Ov*Ov,qc=fp+XM*Ov+fp*Ov*Ov+Ov*Ov*Ov,mn=rd+MD*Ov+MD*Ov*Ov+Ov*Ov*Ov,dM=N6+Ov+ZD*Ov*Ov,zf=gF+rd*Ov+xf*Ov*Ov,YP=XM+MD*Ov+gF*Ov*Ov+Ov*Ov*Ov,DB=fp+MD*Ov+ZD*Ov*Ov+Ov*Ov*Ov,b8=ZD+Ov+ZD*Ov*Ov,JB=gF+gF*Ov+ZD*Ov*Ov+Ov*Ov*Ov,Nv=gF+xf*Ov+Ov*Ov,jS=XM+fp*Ov+fp*Ov*Ov,Nl=rt+rt*Ov+gF*Ov*Ov+Ov*Ov*Ov,Z=xf+ZD*Ov,BP=N6+MD*Ov+XM*Ov*Ov+Ov*Ov*Ov,zb=rt+r1*Ov+ZD*Ov*Ov+Ov*Ov*Ov,WF=gF+Ov+gF*Ov*Ov,YU=N6+XM*Ov+XM*Ov*Ov+Ov*Ov*Ov,Vd=XM+rd*Ov,ID=rd+r1*Ov+Ov*Ov,D3=XM+XM*Ov+MD*Ov*Ov+Ov*Ov*Ov,Vf=N6+rd*Ov+rd*Ov*Ov,B6=rd+Ov+ZD*Ov*Ov,XX=N6+fp*Ov+gF*Ov*Ov+Ov*Ov*Ov,PG=xf+xf*Ov+ZD*Ov*Ov+Ov*Ov*Ov,D=r1+fp*Ov+ZD*Ov*Ov,ZG=xf+rd*Ov+ZD*Ov*Ov+Ov*Ov*Ov,Hl=N6+Ov+ZD*Ov*Ov+Ov*Ov*Ov,j9=xf+gF*Ov+gF*Ov*Ov+Ov*Ov*Ov,Wg=r1+xf*Ov+Ov*Ov,DX=fp+rt*Ov+ZD*Ov*Ov+Ov*Ov*Ov,Bx=r1+gF*Ov+fp*Ov*Ov+Ov*Ov*Ov,jf=N6+rt*Ov+gF*Ov*Ov,pn=N6+fp*Ov+r1*Ov*Ov+Ov*Ov*Ov,Xt=gF+gF*Ov+xf*Ov*Ov,IB=MD+ZD*Ov+gF*Ov*Ov+Ov*Ov*Ov,nb=ZD+Ov+Ov*Ov+Ov*Ov*Ov,nX=fp+ZD*Ov+ZD*Ov*Ov+Ov*Ov*Ov,Y2=rt+XM*Ov+MD*Ov*Ov,mj=MD+MD*Ov+rd*Ov*Ov+Ov*Ov*Ov,Ub=ZD+xf*Ov+XM*Ov*Ov+Ov*Ov*Ov,mb=xf+xf*Ov+Ov*Ov+Ov*Ov*Ov,cX=rd+ZD*Ov+MD*Ov*Ov+Ov*Ov*Ov,l3=XM+XM*Ov+gF*Ov*Ov+Ov*Ov*Ov,lc=fp+MD*Ov+gF*Ov*Ov,pS=ZD+r1*Ov+r1*Ov*Ov,zT=fp+ZD*Ov+r1*Ov*Ov+Ov*Ov*Ov,PS=rt+xf*Ov+gF*Ov*Ov,UB=xf+rt*Ov+ZD*Ov*Ov,NX=fp+Ov+ZD*Ov*Ov,tU=N6+rt*Ov+Ov*Ov+Ov*Ov*Ov,K6=rt+rt*Ov+r1*Ov*Ov,L9=gF+rd*Ov+MD*Ov*Ov+Ov*Ov*Ov,dU=ZD+xf*Ov+Ov*Ov+Ov*Ov*Ov,kX=rd+rd*Ov+Ov*Ov+Ov*Ov*Ov,sf=N6+fp*Ov+Ov*Ov,ND=rt+r1*Ov+ZD*Ov*Ov,kl=gF+MD*Ov+Ov*Ov+Ov*Ov*Ov,hG=ZD+ZD*Ov+ZD*Ov*Ov+Ov*Ov*Ov,Kb=xf+xf*Ov+fp*Ov*Ov+Ov*Ov*Ov,vC=ZD+ZD*Ov+MD*Ov*Ov,F7=ZD+XM*Ov+Ov*Ov+Ov*Ov*Ov,fP=r1+xf*Ov+ZD*Ov*Ov+Ov*Ov*Ov,Gt=gF+gF*Ov,lR=xf+Ov+r1*Ov*Ov+Ov*Ov*Ov,xn=r1+gF*Ov+ZD*Ov*Ov+Ov*Ov*Ov,UM=MD+fp*Ov+r1*Ov*Ov,ff=XM+MD*Ov,kj=xf+Ov+Ov*Ov+Ov*Ov*Ov,f8=fp+rd*Ov+ZD*Ov*Ov,Xf=XM+r1*Ov+Ov*Ov,nx=MD+XM*Ov+XM*Ov*Ov+Ov*Ov*Ov,K2=xf+gF*Ov+MD*Ov*Ov,XB=gF+rd*Ov+rd*Ov*Ov+Ov*Ov*Ov,wl=ZD+rd*Ov+gF*Ov*Ov,Bt=ZD+gF*Ov+fp*Ov*Ov,gr=MD+Ov+gF*Ov*Ov,Yl=XM+ZD*Ov+ZD*Ov*Ov+Ov*Ov*Ov,Jj=rt+rd*Ov+r1*Ov*Ov+Ov*Ov*Ov,gl=xf+r1*Ov+Ov*Ov+Ov*Ov*Ov,p8=N6+Ov+rd*"
```

### `GET /akam/13/32494caf`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"(function(){var _=[\"\\x4d\\x69\\x63\\x72\\x6f\\x73\\x6f\\x66\\x74\\x2e\\x58\\x4d\\x4c\\x48\\x54\\x54\\x50\",\"\\x63\\x72\\x65\\x61\\x74\\x65\\x50\\x6f\\x70\\x75\\x70\",\"\\x5b\",\"\\x73\\x65\\x74\\x41\\x74\\x74\\x72\\x69\\x62\\x75\\x74\\x65\",\"\\x4a\\x53\\x4f\\x4e\",\"\\x6e\\x61\\x70\",\"\\x67\\x79\\x72\\x6f\\x73\\x63\\x6f\\x70\\x65\",\"\\x73\\x75\\x62\\x73\\x74\\x72\",\"\\x64\\x6f\\x4e\\x6f\\x74\\x54\\x72\\x61\\x63\\x6b\",\"\\x73\\x75\\x66\\x66\\x69\\x78\\x65\\x73\",\"\\x7a\\x68\",\"\\x6e\\x61\\x76\",\"\\x61\\x74\\x74\\x61\\x63\\x68\\x45\\x76\\x65\\x6e\\x74\",\"\\x6f\\x72\\x69\\x67\\x69\\x6e\",\"\\x5c\\x0A\",\"\\x63\\x61\\x6d\\x65\\x72\\x61\",\"\\x41\\x42\\x43\\x44\\x45\\x46\\x47\\x48\\x49\\x4a\\x4b\\x4c\\x4d\\x4e\\x4f\\x50\\x51\\x52\\x53\\x54\\x55\\x56\\x57\\x58\\x59\\x5a\\x61\\x62\\x63\\x64\\x65\\x66\\x67\\x68\\x69\\x6a\\x6b\\x6c\\x6d\\x6e\\x6f\\x70\\x71\\x72\\x73\\x74\\x75\\x76\\x77\\x78\\x79\\x7a\\x30\\x31\\x32\\x33\\x34\\x35\\x36\\x37\\x38\\x39\\x2b\\x2f\\x3d\",\"\\x70\\x6f\\x77\",\"\\x6f\\x6e\\x6c\\x6f\\x61\\x64\",\"\\x66\\x69\\x6c\\x6c\\x54\\x65\\x78\\x74\",\"\\x75\\x6e\\x64\\x65\\x66\\x69\\x6e\\x65\\x64\",\"\\x6d\\x73\\x44\\x6f\\x4e\\x6f\\x74\\x54\\x72\\x61\\x63\\x6b\",\"\\x53\\x65\\x71\\x75\\x65\\x6e\\x74\\x75\\x6d\",\"\\x69\\x73\\x20\\x6e\\x6f\\x74\\x20\\x61\\x20\\x76\\x61\\x6c\\x69\\x64\\x20\\x65\\x6e\\x75\\x6d\\x20\\x76\\x61\\x6c\\x75\\x65\\x20\\x6f\\x66\\x20\\x74\\x79\\x70\\x65\\x20\\x50\\x65\\x72\\x6d\\x69\\x73\\x73\\x69\\x6f\\x6e\\x4e\\x61\\x6d\\x65\",\"\\x65\\x78\\x74\\x65\\x72\\x6e\\x61\\x6c\",\"\\x63\\x6c\\x69\\x70\\x62\\x6f\\x61\\x72\\x64\\x2d\\x72\\x65\\x61\\x64\",\"\\x62\\x6f\\x64\\x79\",\"\\x61\\x70\\x70\\x4d\\x69\\x6e\\x6f\\x72\\x56\\x65\\x72\\x73\\x69\\x6f\\x6e\",\"\\x46\\x69\\x72\\x65\\x66\\x6f\\x78\",\"\\x6e\\x75\\x6c\\x6c\",\"\\x49\\x6d\\x61\\x67\\x65\",\"\\x45\\x64\\x67\\x65\",\"\\x5c\\x5c\\x66\",\"\\x74\\x68\\x65\\x6e\",\"\",\"\\x75\\x74\\x66\\x38\\x44\\x65\\x63\\x6f\\x64\\x65\",\"\\x64\\x72\\x69\\x76\\x65\\x72\",\"\\x62\\x65\\x68\\x61\\x76\\x69\\x6f\\x72\",\"\\x7a\",\"\\x53\\x68\\x6f\\x63\\x6b\\x77\\x61\\x76\\x65\\x20\\x46\\x6c\\x61\\x73\\x68\",\"\\x34\",\"\\x61\\x70\\x70\\x6c\\x69\\x63\\x61\\x74\\x69\\x6f\\x6e\\x2f\\x78\\x2d\\x73\\x68\\x6f\\x63\\x6b\\x77\\x61\\x76\\x65\\x2d\\x66\\x6c\\x61\\x73\\x68\",\"\\x6c\\x61\\x6e\\x67\\x75\\x61\\x67\\x65\\x73\",\"\\x73\\x65\\x6c\\x65\\x6e\\x69\\x75\\x6d\",\"\\x5b\\x6f\\x62\\x6a\\x65\\x63\\x74\\x20\\x41\\x72\\x72\\x61\\x79\\x5d\",\"\\x50\\x4f\\x53\\x54\",\"\\x74\\x79\\x70\\x65\",\"\\x37\",\"\\x61\\x63\\x63\\x65\\x6c\\x65\\x72\\x6f\\x6d\\x65\\x74\\x65\\x72\",\"\\x53\\x69\\x6c\\x76\\x65\\x72\\x6c\\x69\\x67\\x68\\x74\\x20\\x50\\x6c\\x75\\x67\\x2d\\x49\\x6e\",\"\\x64\\x65\\x76\\x69\\x63\\x65\\x2d\\x69\\x6e\\x66\\x6f\",\"\\x74\\x6f\\x44\\x61\\x74\\x61\\x55\\x52\\x4c\",\"\\x73\\x72\",\"\\x72\\x65\\x73\\x70\\x6f\\x6e\\x73\\x65\\x53\\x74\\x61\\x72\\x74\",\"\\x73\\x70\",\"\\x70\\x72\\x6f\\x74\\x6f\\x74\\x79\\x70\\x65\",\"\\x64\\x72\\x61\\x77\\x49\\x6d\\x61\\x67\\x65\",\"\\x32\\x2e\\x30\",\"\\x67\\x65\\x74\\x41\\x74\\x74\\x72\\x69\\x62\\x75\\x74\\x65\",\"\\x5f\\x5f\\x61\\x6b\\x66\\x70\\x5f\\x73\\x74\\x6f\\x72\\x61\\x67\\x65\\x5f\\x74\\x65\\x73\\x74\\x5f\\x5f\",\"\\x6d\\x69\\x63\\x72\\x6f\\x70\\x68\\x6f\\x6e\\x65\",\"\\x5c\\x5c\\x75\",\"\\x64\\x61\\x74\\x61\",\"\\x73\\x70\\x65\\x61\\x6b\\x65\\x72\",\"\\x53\\x6f\\x66\\x74\\x20\\x52\\x75\\x64\\x64\\x79\\x20\\x46\\x6f\\x6f\\x74\\x68\\x6f\\x6c\\x64\\x20\\x32\",\"\\x52\\x4f\\x54\\x4c\",\"\\x66\\x73\\x66\\x70\",\"\\x6d\\x61\\x70\",\"\\x3a\\x20\",\"\\x63\\x6f\\x6e\\x74\\x65\\x78\\x74\\x4d\\x65\\x6e\\x75\",\"\\x6c\\x61\\x6e\\x67\\x75\\x61\\x67\\x65\",\"\\x58\\x50\\x61\\x74\\x68\\x52\\x65\\x73\\x75\\x6c\\x74\",\"\\x73\\x74\\x72\\x69\\x6e\\x67\\x69\\x66\\x79\",\"\\x50\\x44\\x46\\x2e\\x50\\x64\\x66\\x43\\x74\\x72\\x6c\\x2e\\x31\",\"\\x63\\x72\\x65\\x61\\x74\\x65\\x45\\x6c\\x65\\x6d\\x65\\x6e\\x74\",\"\\x73\\x61\\x76\\x65\",\"\\x69\\x6d\\x61\\x67\\x65\\x73\",\"\\x72\\x67\\x62\\x61\\x28\\x32\\x35\\x35\\x2c\\x31\\x35\\x33\\x2c\\x31\\x35\\x33\\x2c\\x20\\x30\\x2e\\x35\\x29\",\"\\x62\\x61\\x7a\\x61\\x64\\x65\\x62\\x65\\x7a\\x6f\\x6c\\x6b\\x6f\\x68\\x70\\x65\\x70\\x61\\x64\\x72\",\"\\x63\\x72\\x63\",\"\\x31\\x2e\\x32\",\"\\x63\\x68\\x61\\x72\\x43\\x6f\\x64\\x65\\x41\\x74\",\"\\x61\\x70\\x70\\x4e\\x61\\x6d\\x65\",\"\\x70\\x65\\x72\\x66\\x6f\\x72\\x6d\\x61\\x6e\\x63\\x65\",\"\\x74\\x6f\\x4a\\x53\\x4f\\x4e\",\"\\x63\\x76\",\"\\x41\\x63\\x74\\x69\\x76\\x65\\x58\\x4f\\x62\\x6a\\x65\\x63\\x74\",\"\\x75\\x72\\x6c\\x28\\x23\\x64\\x65\\x66\\x61\\x75\\x6c\\x74\\x23\\x75\\x73\\x65\\x72\\x44\\x61\\x74\\x61\\x29\",\"\\x70\\x75\\x73\\x68\",\"\\x67\\x65\\x74\\x45\\x6c\\x65\\x6d\\x65\\x6e\\x74\\x73\\x42\\x79\\x54\\x61\\x67\\x4e\\x61\\x6d\\x65\",\"\\x68\\x6f\\x73\\x74\\x6e\\x61\\x6d\\x65\",\"\\x4a\\x61\\x76\\x61\\x53\\x63\\x72\\x69\\x70\\x74\",\"\\x6d\\x61\\x74\\x63\\x68\",\"\\x31\\x2e\\x30\",\"\\x61\\x70\\x70\\x56\\x65\\x72\\x73\\x69\\x6f\\x6e\",\"\\x68\\x61\\x73\\x4f\\x77\\x6e\\x50\\x72\\x6f\\x70\\x65\\x72\\x74\\x79\",\"\\x72\\x
```

### `GET /areanaologada/2.0.0/lista-conteudo-servico`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "distribuidora-cms": "5",
  "distribuidora": "NEOENERGIA",
  "service_name": "home-anl"
}
```
response body:
```json
[
  {
    "contentId": 61,
    "contentName": "destaque-anl-link-03",
    "attributeName": "Texto",
    "contentVersion": {
      "contentVersionId": 34,
      "url": "/content/text/destaque-anl-link-03",
      "message": "https://www.neoenergia.com/economia-circular",
      "dtInicial": "20-05-2025 16:51:35",
      "dtFinal": null,
      "file": null
    }
  },
  {
    "contentId": 57,
    "contentName": "destaque-anl-link-02",
    "attributeName": "Texto",
    "contentVersion": {
      "contentVersionId": 33,
      "url": "/content/text/destaque-anl-link-02",
      "message": "https://www.neoenergia.com/seguranca-na-internet",
      "dtInicial": "20-05-2025 16:51:08",
      "dtFinal": null,
      "file": null
    }
  },
  {
    "contentId": 49,
    "contentName": "destaque-anl-texto-01",
    "attributeName": "Texto",
    "contentVersion": {
      "contentVersionId": 23,
      "url": "/content/text/conteudo-destaque1",
      "message": "<span style=\"font-size:12.0pt;line-height:115%;\r\nfont-family:&quot;Aptos&quot;,sans-serif;mso-ascii-theme-font:minor-latin;mso-fareast-font-family:\r\nAptos;mso-fareast-theme-font:minor-latin;mso-hansi-theme-font:minor-latin;\r\nmso-bidi-font-family:&quot;Times New Roman&quot;;mso-bidi-theme-font:minor-bidi;\r\nmso-ansi-language:PT-BR;mso-fareast-language:EN-US;mso-bidi-language:AR-SA\">Cadastre-se\r\ne tenha acesso a todos os serviços da Agência Virtual de maneira simples pelo\r\nnosso site. Veja como essa novidade pode facilitar os seus pagamentos e\r\ncontrole do seu consumo mensal.</span>",
      "dtInicial": "19-05-2025 18:17:17",
      "dtFinal": null,
      "file": null
    }
  },
  {
    "contentId": 47,
    "contentName": "destaque-anl-titulo-01",
    "attributeName": "Texto",
    "contentVersion": {
      "contentVersionId": 21,
      "url": "/content/text/titulo-destaque1",
      "message": "Novo login, como fazer?",
      "dtInicial": "15-05-2025 11:14:36",
      "dtFinal": null,
      "file": null
    }
  },
  {
    "contentId": 53,
    "contentName": "banner-anl-titulo-01",
    "attributeName": "Texto",
    "contentVersion": {
      "contentVersionId": 31,
      "url": "/content/text/banner-anl-titulo-01",
      "message": "Mais por você",
      "dtInicial": "20-05-2025 01:02:00",
      "dtFinal": null,
      "file": null
    }
  },
  {
    "contentId": 55,
    "contentName": "destaque-anl-titulo-02",
    "attributeName": "Texto",
    "contentVersion": {
      "contentVersionId": 25,
      "url": "/content/text/destaque-anl-titulo-02",
      "message": "Não caia em golpes na internet",
      "dtInicial": "20-05-2025 03:58:00",
      "dtFinal": null,
      "file": null
    }
  },
  {
    "contentId": 56,
    "contentName": "destaque-anl-texto-02",
    "attributeName": "Texto",
    "contentVersion": {
      "contentVersionId": 26,
      "url": "/content/text/destaque-anl-texto-02",
      "message": "<p class=\"MsoNormal\">Conheça os principais golpes na Internet e aprenda cinco\r\ndicas práticas para não cair em fraudes digitais.<o:p></o:p></p>",
      "dtInicial": "19-05-2025 18:58:52",
      "dtFinal": null,
      "file": null
    }
  },
  {
    "contentId": 58,
    "contentName": "destaque-anl-banner-03",
    "attributeName": "Imagem",
    "contentVersion": {
      "contentVersionId": 28,
      "url": "/content/image/destaque-anl-banner-03",
      "message": "",
      "dtInicial": "19-05-2025 19:00:13",
      "dtFinal": null,
      "file": "/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAUEBAUEAwUFBAUGBgUGCA4JCAcHCBEMDQoOFBEVFBMRExMWGB8bFhceFxMTGyUcHiAhIyMjFRomKSYiKR8iIyL/2wBDAQYGBggHCBAJCRAiFhMWIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiL/wgARCALYBLADAREAAhEBAxEB/8QAHQAAAQUBAQEBAAAAAAAAAAAAAwECBAUGAAcICf/EABoBAAMBAQEBAAAAAAAAAAAAAAABAgMEBQb/2gAMAwEAAhADEAAAAPRshQQai4fC5nI4FGjXJo1wIDRINAahrGgiGg0GiaxAYJoxNBCIzG3GOuNJnpu8NrPO4NTQ7YRqI4o7lkXKw6ZeHQbn3lc/RdY6Wyk4laLcG0lakbTGhNDYypGJrGORtLaNcFpMaiiq6K4UcChf51fhPCeOycy6D0FacxwlDgUFDgUOBQ4OBQ4OFw+BA4XD4ABRZXm8NWiNRNpSRPY8EBoICoaOPLHLcMcuMK/3yn3L2jCqMHXc2llasdEdBUyyy0iWrjpi22k7OAaEZyGg4CsQFaYmCTgkUPZ8+Zrg4ahwcHC4ODmcmgcJGICA0GggNBoNQ1jRNBEMY0GAMANRRZmjM6RcQ9zz73OW1fU5Pr5RAEQhTuPut+Hum8uws9QzpMJsCZTRUEachGNBoMoYkOhoNaa0y0+4kaRI0gbUcIAQGorRam1Jt3M5qwHZOZVEho1IjTmlBzFQocCgocHAocHBwcCBwcJR8KOOhyvNZaKiYydSI01DwcxARPhNGxMEvkxpxktBvlLqQyCgZNMzocqZSIER3Em4VFO2m0ytouNEeUmdVkPorQJaBxqOrKfak0nAjShw1D57zXAo+Dg4XD4FGguBA4EBGmggNBAaNgmg0GiaJoNGwGMGgDkIVTda0UWjx3usda65yfTysTE0IJnN12PH2dnddcw9ZBeTglw7vK7aLdNNhsTRUipJfTSJqn0vm1Eak+ksjW2IA2Np1SaolaZFuZlxZuJ9KVcS9ImaQdojHCcCsVHAocChwcHAocCBwICgggDocrzeOpWprJtSQBJuBzSj4TR8hqAqmICqiov9sTNUXNULMstXJbhQ0gqafmHS8ntES5jOWEdSLIwAUotICmTN7bm6fVObX0lz6H3c1jpBWOEg2B4BBwuGocHC4fC4fBwILgRjQQEBoNBgNQjGg0TQaDGMQIQxiaE3HJijETcY7XWW1bU5Lp5hAOpGm1DVcdpgLNBpCcGT1HN0aXHWQmNAVUcdUVR1o5ElTYyTpOlon01ypYaKuVNVKjmCGFONaVzayrUU7TK86ua+6cLHSDg9pwKCgocHAocHAocHBwIHBwcEdFDnecw0kUWFKU04EBwOBAaNqEQoBTFLjqoiL7fKBL8/5dKVufqpjYJK6zzfrww/RgoGQVBQemgNaVMbgLQRKFlnrqefo9Yl+2dEb28pDTUIHgUiggKHAglGoIjmILmINBIDQ
```

### `GET /assets/fonts//IberPangea/IberPangeaText-Bold.ttf`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<unreadable: 'utf-8' codec can't decode byte 0xdb in position 22: invalid continuation byte>"
```

### `GET /assets/fonts//IberPangea/IberPangeaText-Regular.ttf`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<unreadable: 'utf-8' codec can't decode byte 0xf2 in position 22: invalid continuation byte>"
```

### `GET /assets/images/icons/account_circle.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"20\" height=\"21\" viewBox=\"0 0 20 21\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<path d=\"M3.85 15.1117C4.7 14.4617 5.65 13.9492 6.7 13.5742C7.75 13.1992 8.85 13.0117 10 13.0117C11.15 13.0117 12.25 13.1992 13.3 13.5742C14.35 13.9492 15.3 14.4617 16.15 15.1117C16.7333 14.4284 17.1875 13.6534 17.5125 12.7867C17.8375 11.9201 18 10.9951 18 10.0117C18 7.79505 17.2208 5.90755 15.6625 4.34922C14.1042 2.79089 12.2167 2.01172 10 2.01172C7.78333 2.01172 5.89583 2.79089 4.3375 4.34922C2.77917 5.90755 2 7.79505 2 10.0117C2 10.9951 2.1625 11.9201 2.4875 12.7867C2.8125 13.6534 3.26667 14.4284 3.85 15.1117ZM10 11.0117C9.01667 11.0117 8.1875 10.6742 7.5125 9.99922C6.8375 9.32422 6.5 8.49505 6.5 7.51172C6.5 6.52839 6.8375 5.69922 7.5125 5.02422C8.1875 4.34922 9.01667 4.01172 10 4.01172C10.9833 4.01172 11.8125 4.34922 12.4875 5.02422C13.1625 5.69922 13.5 6.52839 13.5 7.51172C13.5 8.49505 13.1625 9.32422 12.4875 9.99922C11.8125 10.6742 10.9833 11.0117 10 11.0117ZM10 20.0117C8.61667 20.0117 7.31667 19.7492 6.1 19.2242C4.88333 18.6992 3.825 17.9867 2.925 17.0867C2.025 16.1867 1.3125 15.1284 0.7875 13.9117C0.2625 12.6951 0 11.3951 0 10.0117C0 8.62839 0.2625 7.32839 0.7875 6.11172C1.3125 4.89505 2.025 3.83672 2.925 2.93672C3.825 2.03672 4.88333 1.32422 6.1 0.799219C7.31667 0.274219 8.61667 0.0117188 10 0.0117188C11.3833 0.0117188 12.6833 0.274219 13.9 0.799219C15.1167 1.32422 16.175 2.03672 17.075 2.93672C17.975 3.83672 18.6875 4.89505 19.2125 6.11172C19.7375 7.32839 20 8.62839 20 10.0117C20 11.3951 19.7375 12.6951 19.2125 13.9117C18.6875 15.1284 17.975 16.1867 17.075 17.0867C16.175 17.9867 15.1167 18.6992 13.9 19.2242C12.6833 19.7492 11.3833 20.0117 10 20.0117ZM10 18.0117C10.8833 18.0117 11.7167 17.8826 12.5 17.6242C13.2833 17.3659 14 16.9951 14.65 16.5117C14 16.0284 13.2833 15.6576 12.5 15.3992C11.7167 15.1409 10.8833 15.0117 10 15.0117C9.11667 15.0117 8.28333 15.1409 7.5 15.3992C6.71667 15.6576 6 16.0284 5.35 16.5117C6 16.9951 6.71667 17.3659 7.5 17.6242C8.28333 17.8826 9.11667 18.0117 10 18.0117ZM10 9.01172C10.4333 9.01172 10.7917 8.87005 11.075 8.58672C11.3583 8.30338 11.5 7.94505 11.5 7.51172C11.5 7.07839 11.3583 6.72005 11.075 6.43672C10.7917 6.15339 10.4333 6.01172 10 6.01172C9.56667 6.01172 9.20833 6.15339 8.925 6.43672C8.64167 6.72005 8.5 7.07839 8.5 7.51172C8.5 7.94505 8.64167 8.30338 8.925 8.58672C9.20833 8.87005 9.56667 9.01172 10 9.01172Z\" fill=\"#00402A\"/>\n</svg>\n"
```

### `GET /assets/images/icons/arrow_forward.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"11\" height=\"21\" viewBox=\"0 0 11 21\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<path d=\"M1.65817 20.594L0 18.8014L7.68365 10.495L0 2.18857L1.65817 0.395996L11 10.495L1.65817 20.594Z\" fill=\"#393735\"/>\n</svg>\n"
```

### `GET /assets/images/icons/atencao.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"213.016\" height=\"190.05\" viewBox=\"0 0 213.016 190.05\">\n  <path id=\"warning_amber-24px\" d=\"M108.977,26.313l84.144,145.38H24.833l84.144-145.38M5.5,160.519a22.344,22.344,0,0,0,19.332,33.523H193.12a22.344,22.344,0,0,0,19.332-33.523L128.308,15.139a22.337,22.337,0,0,0-38.664,0ZM97.8,82.3v22.349a11.174,11.174,0,0,0,22.349,0V82.3a11.174,11.174,0,0,0-22.349,0Zm0,55.872h22.349v22.349H97.8Z\" transform=\"translate(-2.469 -3.992)\" fill=\"#D60E00\"/>\n</svg>\n"
```

### `GET /assets/images/icons/facebook_silver_v2.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"20\" height=\"21\" viewBox=\"0 0 20 21\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<g clip-path=\"url(#clip0_2163_2018)\">\n<g clip-path=\"url(#clip1_2163_2018)\">\n<path d=\"M10 20.21C7.69141 20.2116 5.45335 19.4143 3.66566 17.9535C1.87796 16.4927 0.650769 14.4584 0.192328 12.1958C-0.266113 9.93312 0.0724422 7.58155 1.15054 5.54011C2.22864 3.49868 3.97986 1.89315 6.10703 0.995987C8.55073 -0.0363683 11.3045 -0.0556829 13.7624 0.942292C16.2203 1.94027 18.1812 3.87378 19.2135 6.31749C20.2459 8.76119 20.2652 11.5149 19.2672 13.9729C18.2692 16.4308 16.3357 18.3916 13.892 19.424C12.6608 19.9451 11.337 20.2125 10 20.21ZM7.66703 8.77499V10.65H8.79103V16.09H11.042V10.65H12.542L12.742 8.77499H11.042V7.83699C11.042 7.35699 11.077 7.08599 11.791 7.08599H12.73V5.20999H11.23C10.8856 5.17693 10.5382 5.21742 10.2106 5.32879C9.88305 5.44016 9.5829 5.61986 9.33003 5.85599C8.92836 6.36325 8.73428 7.00411 8.78703 7.64899V8.77499H7.66703Z\" fill=\"#007F33\"/>\n</g>\n</g>\n<defs>\n<clipPath id=\"clip0_2163_2018\">\n<rect width=\"20\" height=\"20\" fill=\"white\" transform=\"translate(0 0.209961)\"/>\n</clipPath>\n<clipPath id=\"clip1_2163_2018\">\n<rect width=\"20\" height=\"20\" fill=\"white\" transform=\"translate(0 0.209961)\"/>\n</clipPath>\n</defs>\n</svg>\n"
```

### `GET /assets/images/icons/formulario_cadastro.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"32.922\" height=\"36.58\" viewBox=\"0 0 32.922 36.58\">\n  <path id=\"ic_assignment_24px\" d=\"M32.264,4.658H24.619a5.465,5.465,0,0,0-10.316,0H6.658A3.669,3.669,0,0,0,3,8.316V33.922A3.669,3.669,0,0,0,6.658,37.58H32.264a3.669,3.669,0,0,0,3.658-3.658V8.316A3.669,3.669,0,0,0,32.264,4.658Zm-12.8,0a1.829,1.829,0,1,1-1.829,1.829A1.834,1.834,0,0,1,19.461,4.658Zm3.658,25.606h-12.8V26.606h12.8Zm5.487-7.316H10.316V19.29h18.29Zm0-7.316H10.316V11.974h18.29Z\" transform=\"translate(-3 -1)\" fill=\"#00A443\"/>\n</svg>\n"
```

### `GET /assets/images/icons/icone-vincular-uc.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<mask id=\"0_1708_61389\" style=\"mask-type:alpha\" maskUnits=\"userSpaceOnUse\" x=\"0\" y=\"0\" width=\"24\" height=\"24\">\n<rect width=\"24\" height=\"24\" fill=\"#D9D9D9\"/>\n</mask>\n<g mask=\"url(#mask0_1708_61389)\">\n<path d=\"M23 5V19C23 19.55 22.8042 20.0208 22.4125 20.4125C22.0208 20.8042 21.55 21 21 21H18C17.7167 21 17.4792 20.9042 17.2875 20.7125C17.0958 20.5208 17 20.2833 17 20C17 19.7167 17.0958 19.4792 17.2875 19.2875C17.4792 19.0958 17.7167 19 18 19H21V5H12V5.5C12 5.78333 11.9042 6.02083 11.7125 6.2125C11.5208 6.40417 11.2833 6.5 11 6.5C10.7167 6.5 10.4792 6.40417 10.2875 6.2125C10.0958 6.02083 10 5.78333 10 5.5V4.95C10 4.41667 10.1917 3.95833 10.575 3.575C10.9583 3.19167 11.4167 3 11.95 3H21C21.55 3 22.0208 3.19583 22.4125 3.5875C22.8042 3.97917 23 4.45 23 5ZM1 12.025C1 11.6917 1.075 11.3833 1.225 11.1C1.375 10.8167 1.58333 10.5833 1.85 10.4L6.85 6.825C7.03333 6.69167 7.22083 6.59583 7.4125 6.5375C7.60417 6.47917 7.8 6.45 8 6.45C8.2 6.45 8.39583 6.47917 8.5875 6.5375C8.77917 6.59583 8.96667 6.69167 9.15 6.825L14.15 10.4C14.4167 10.5833 14.625 10.8167 14.775 11.1C14.925 11.3833 15 11.6917 15 12.025V19C15 19.55 14.8042 20.0208 14.4125 20.4125C14.0208 20.8042 13.55 21 13 21H11C10.45 21 9.97917 20.8042 9.5875 20.4125C9.19583 20.0208 9 19.55 9 19V16H7V19C7 19.55 6.80417 20.0208 6.4125 20.4125C6.02083 20.8042 5.55 21 5 21H3C2.45 21 1.97917 20.8042 1.5875 20.4125C1.19583 20.0208 1 19.55 1 19V12.025ZM3 12V19H5V16C5 15.45 5.19583 14.9792 5.5875 14.5875C5.97917 14.1958 6.45 14 7 14H9C9.55 14 10.0208 14.1958 10.4125 14.5875C10.8042 14.9792 11 15.45 11 16V19H13V12L8 8.45L3 12ZM17.5 9H18.5C18.6333 9 18.75 8.95 18.85 8.85C18.95 8.75 19 8.63333 19 8.5V7.5C19 7.36667 18.95 7.25 18.85 7.15C18.75 7.05 18.6333 7 18.5 7H17.5C17.3667 7 17.25 7.05 17.15 7.15C17.05 7.25 17 7.36667 17 7.5V8.5C17 8.63333 17.05 8.75 17.15 8.85C17.25 8.95 17.3667 9 17.5 9ZM17.5 13H18.5C18.6333 13 18.75 12.95 18.85 12.85C18.95 12.75 19 12.6333 19 12.5V11.5C19 11.3667 18.95 11.25 18.85 11.15C18.75 11.05 18.6333 11 18.5 11H17.5C17.3667 11 17.25 11.05 17.15 11.15C17.05 11.25 17 11.3667 17 11.5V12.5C17 12.6333 17.05 12.75 17.15 12.85C17.25 12.95 17.3667 13 17.5 13ZM17.5 17H18.5C18.6333 17 18.75 16.95 18.85 16.85C18.95 16.75 19 16.6333 19 16.5V15.5C19 15.3667 18.95 15.25 18.85 15.15C18.75 15.05 18.6333 15 18.5 15H17.5C17.3667 15 17.25 15.05 17.15 15.15C17.05 15.25 17 15.3667 17 15.5V16.5C17 16.6333 17.05 16.75 17.15 16.85C17.25 16.95 17.3667 17 17.5 17Z\" fill=\"#007F33\"/>\n</g>\n</svg>\n"
```

### `GET /assets/images/icons/icone_fatura_facil.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg id=\"copy\" xmlns=\"http://www.w3.org/2000/svg\" width=\"33.74\" height=\"41.238\" viewBox=\"0 0 33.74 41.238\">\n  <path id=\"Caminho_332\" data-name=\"Caminho 332\" d=\"M12.754,6.918H7V5h5.754Z\" transform=\"translate(2.329 2.864)\" fill=\"#fff\"/>\n  <path id=\"Caminho_333\" data-name=\"Caminho 333\" d=\"M12.754,10.918H7V9h5.754Z\" transform=\"translate(2.329 6.728)\" fill=\"#fff\"/>\n  <path id=\"Caminho_334\" data-name=\"Caminho 334\" d=\"M7,14.918h5.754V13H7Z\" transform=\"translate(2.329 10.592)\" fill=\"#fff\"/>\n  <path id=\"Caminho_335\" data-name=\"Caminho 335\" d=\"M3,34.74V1H29.242V8.5h7.5v33.74H10.5v-7.5Zm22.494-3.749V4.749H6.749V30.991Zm3.749-18.745V34.74h-15v3.749H32.991V12.247Z\" transform=\"translate(-3 -1)\" fill=\"#fff\" fill-rule=\"evenodd\"/>\n</svg>\n"
```

### `GET /assets/images/icons/icone_ligacao_nova.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24.084\" height=\"41.239\" viewBox=\"0 0 24.084 41.239\">\n  <path id=\"Caminho_14530\" data-name=\"Caminho 14530\" d=\"M733.275,84.458h-11.2l7.939-16.675h-10.7L709.191,92.346h9.161l-8.106,16.676Z\" transform=\"translate(-709.191 -67.783)\" fill=\"#00A443\"/>\n</svg>\n"
```

### `GET /assets/images/icons/instagram_silver_v2.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"20\" height=\"21\" viewBox=\"0 0 20 21\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<g clip-path=\"url(#clip0_2163_2029)\">\n<g clip-path=\"url(#clip1_2163_2029)\">\n<path d=\"M10 20.21C7.69141 20.2116 5.45335 19.4143 3.66566 17.9535C1.87796 16.4927 0.650769 14.4584 0.192328 12.1958C-0.266113 9.93312 0.0724422 7.58155 1.15054 5.54011C2.22864 3.49868 3.97986 1.89315 6.10703 0.995987C8.55073 -0.0363683 11.3045 -0.0556829 13.7624 0.942292C16.2203 1.94027 18.1812 3.87378 19.2135 6.31749C20.2459 8.76119 20.2652 11.5149 19.2672 13.9729C18.2692 16.4308 16.3357 18.3916 13.892 19.424C12.6608 19.9451 11.337 20.2125 10 20.21ZM10 4.87699H9.80503C8.54703 4.87699 8.35103 4.88399 7.80503 4.90799C7.36283 4.91742 6.92538 5.00126 6.51103 5.15599C6.15665 5.29295 5.83481 5.50247 5.56616 5.77112C5.29752 6.03977 5.08799 6.36161 4.95103 6.71599C4.79491 7.13006 4.71004 7.56758 4.70003 8.00999C4.67403 8.58899 4.66803 8.77999 4.66803 10.21C4.66803 11.64 4.67403 11.836 4.70003 12.41C4.70961 12.8522 4.79345 13.2896 4.94803 13.704C5.08504 14.0582 5.29447 14.3799 5.56292 14.6485C5.83137 14.9171 6.15293 15.1267 6.50703 15.264C6.92122 15.4178 7.3583 15.5009 7.80003 15.51C8.37303 15.536 8.56203 15.542 10 15.542C11.438 15.542 11.626 15.536 12.2 15.51C12.6442 15.5011 13.0838 15.4173 13.5 15.262C13.8542 15.1246 14.1758 14.915 14.4444 14.6464C14.713 14.3778 14.9227 14.0561 15.06 13.702C15.2111 13.2876 15.2922 12.851 15.3 12.41C15.326 11.831 15.332 11.64 15.332 10.21C15.332 8.77999 15.326 8.58999 15.3 8.00999C15.2898 7.56787 15.206 7.13056 15.052 6.71599C14.9149 6.3617 14.7053 6.03995 14.4367 5.77132C14.1681 5.50269 13.8463 5.29311 13.492 5.15599C13.0781 5.00247 12.6414 4.91932 12.2 4.90999C11.625 4.88299 11.435 4.87699 10 4.87699Z\" fill=\"#007F33\"/>\n<path d=\"M10.0011 14.491C8.55411 14.491 8.38211 14.486 7.81111 14.46C7.46756 14.4555 7.12732 14.3922 6.80511 14.273C6.56986 14.1861 6.357 14.0477 6.18211 13.868C6.00196 13.6932 5.86323 13.4803 5.77611 13.245C5.65662 12.9228 5.59337 12.5825 5.58911 12.239C5.56311 11.673 5.55811 11.506 5.55811 10.048C5.55811 8.58996 5.56311 8.42496 5.58911 7.85896C5.59326 7.51505 5.65651 7.17442 5.77611 6.85196C5.86306 6.61674 6.00143 6.4039 6.18111 6.22896C6.35619 6.04874 6.56941 5.91001 6.80511 5.82296C7.12726 5.70347 7.46754 5.64022 7.81111 5.63596C8.31111 5.61296 8.50311 5.60696 9.51111 5.60596H9.99711C11.4581 5.60596 11.6271 5.61296 12.1911 5.63896C12.5349 5.64292 12.8755 5.70584 13.1981 5.82496C13.4335 5.91161 13.6464 6.05002 13.8211 6.22996C14.0009 6.40481 14.1393 6.61768 14.2261 6.85296C14.3456 7.17511 14.4088 7.51539 14.4131 7.85896C14.4391 8.42896 14.4451 8.59996 14.4451 10.049C14.4451 11.498 14.4391 11.668 14.4131 12.238C14.4088 12.5815 14.3456 12.9218 14.2261 13.244C14.1358 13.4776 13.9977 13.6897 13.8207 13.8669C13.6437 14.0441 13.4316 14.1824 13.1981 14.273C12.876 14.3928 12.5357 14.456 12.1921 14.46C11.6201 14.486 11.4491 14.491 10.0011 14.491ZM10.0011 7.26696C9.45072 7.26696 8.9127 7.43015 8.45506 7.73591C7.99742 8.04166 7.64071 8.47625 7.43004 8.98472C7.21938 9.49319 7.16421 10.0527 7.27151 10.5925C7.37881 11.1323 7.64376 11.6282 8.03287 12.0175C8.42198 12.4067 8.91777 12.6719 9.45756 12.7794C9.99734 12.8869 10.5569 12.8319 11.0654 12.6214C11.574 12.4109 12.0087 12.0544 12.3146 11.5968C12.6205 11.1393 12.7839 10.6013 12.7841 10.051C12.7833 9.31301 12.4899 8.6055 11.9682 8.0836C11.4465 7.5617 10.7391 7.26802 10.0011 7.26696ZM12.8951 6.50596C12.7665 6.50596 12.6409 6.54408 12.534 6.6155C12.4271 6.68692 12.3438 6.78844 12.2946 6.90721C12.2454 7.02598 12.2325 7.15668 12.2576 7.28277C12.2827 7.40885 12.3446 7.52467 12.4355 7.61558C12.5264 7.70648 12.6422 7.76839 12.7683 7.79347C12.8944 7.81855 13.0251 7.80568 13.1438 7.75648C13.2626 7.70728 13.3641 7.62397 13.4356 7.51708C13.507 7.41019 13.5451 7.28451 13.5451 7.15596C13.5448 6.98365 13.4763 6.81847 13.3544 6.69663C13.2326 6.57479 13.0674 6.50622 12.8951 6.50596Z\" fill=\"#007F33\"/>\n</g>\n</g>\n<defs>\n<clipPath id=\"clip0_2163_2029\">\n<rect width=\"20\" height=\"20\" fill=\"white\" "
```

### `GET /assets/images/icons/lampada.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"25.893\" height=\"25\" viewBox=\"0 0 25.893 25\">\n  <path id=\"Caminho_14209\" data-name=\"Caminho 14209\" d=\"M7.829,22.727h4.6a2.3,2.3,0,0,1-4.6,0Zm-2.3-2.273a1.147,1.147,0,0,0,1.151,1.136h6.9a1.136,1.136,0,1,0,0-2.273h-6.9A1.147,1.147,0,0,0,5.528,20.455ZM18.762,10.8a8.554,8.554,0,0,1-4.338,7.386H5.838A8.554,8.554,0,0,1,1.5,10.8a8.58,8.58,0,0,1,8.631-8.523A8.58,8.58,0,0,1,18.762,10.8Zm5.6-2.42-1.577.716,1.577.716.725,1.557.725-1.557,1.577-.716-1.577-.716-.725-1.557ZM21.639,6.818l1.082-2.341,2.371-1.068L22.721,2.341,21.639,0,20.557,2.341,18.187,3.409l2.371,1.068Z\" transform=\"translate(-1.5)\" fill=\"#007F33\"/>\n</svg>\n"
```

### `GET /assets/images/icons/linkedin-silver.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"20\" height=\"21\" viewBox=\"0 0 20 21\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<g clip-path=\"url(#clip0_1981_5058)\">\n<path d=\"M10 0.5C4.47656 0.5 0 4.97656 0 10.5C0 16.0234 4.47656 20.5 10 20.5C15.5234 20.5 20 16.0234 20 10.5C20 4.97656 15.5234 0.5 10 0.5ZM7.20312 14.6836H5.25391V8.44531H7.20312V14.6836ZM6.17578 7.66406H6.16016C5.45312 7.66406 4.99609 7.1875 4.99609 6.58203C4.99609 5.96484 5.46875 5.5 6.1875 5.5C6.90625 5.5 7.34766 5.96484 7.36328 6.58203C7.36719 7.18359 6.91016 7.66406 6.17578 7.66406ZM15 14.6836H12.7891V11.457C12.7891 10.6133 12.4453 10.0352 11.6836 10.0352C11.1016 10.0352 10.7773 10.4258 10.6289 10.8008C10.5742 10.9336 10.582 11.1211 10.582 11.3125V14.6836H8.39063C8.39063 14.6836 8.41797 8.96484 8.39063 8.44531H10.582V9.42578C10.7109 8.99609 11.4102 8.38672 12.5273 8.38672C13.9141 8.38672 15 9.28516 15 11.2148V14.6836Z\" fill=\"#007F33\"/>\n</g>\n<defs>\n<clipPath id=\"clip0_1981_5058\">\n<rect width=\"20\" height=\"20\" fill=\"white\" transform=\"translate(0 0.5)\"/>\n</clipPath>\n</defs>\n</svg>\n"
```

### `GET /assets/images/icons/logout.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"18\" height=\"18\" viewBox=\"0 0 18 18\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<path d=\"M2 18C1.45 18 0.979167 17.8042 0.5875 17.4125C0.195833 17.0208 0 16.55 0 16V2C0 1.45 0.195833 0.979167 0.5875 0.5875C0.979167 0.195833 1.45 0 2 0H9V2H2V16H9V18H2ZM13 14L11.625 12.55L14.175 10H6V8H14.175L11.625 5.45L13 4L18 9L13 14Z\" fill=\"#393735\"/>\n</svg>\n"
```

### `GET /assets/images/icons/person.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"16\" height=\"19\" viewBox=\"0 0 16 19\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<path d=\"M6 10.0117C5.16667 10.0117 4.45833 9.72005 3.875 9.13672C3.29167 8.55339 3 7.84505 3 7.01172C3 6.64505 3.05833 6.29922 3.175 5.97422C3.29167 5.64922 3.46667 5.35339 3.7 5.08672C3.63333 4.92005 3.58333 4.74505 3.55 4.56172C3.51667 4.37839 3.5 4.19505 3.5 4.01172C3.5 3.37839 3.67083 2.81589 4.0125 2.32422C4.35417 1.83255 4.8 1.47005 5.35 1.23672C5.68333 0.853385 6.075 0.553385 6.525 0.336719C6.975 0.120052 7.46667 0.0117188 8 0.0117188C8.53333 0.0117188 9.025 0.120052 9.475 0.336719C9.925 0.553385 10.3167 0.853385 10.65 1.23672C11.2 1.47005 11.6458 1.83255 11.9875 2.32422C12.3292 2.81589 12.5 3.37839 12.5 4.01172C12.5 4.19505 12.4833 4.37839 12.45 4.56172C12.4167 4.74505 12.3667 4.92005 12.3 5.08672C12.5333 5.35339 12.7083 5.64922 12.825 5.97422C12.9417 6.29922 13 6.64505 13 7.01172C13 7.84505 12.7083 8.55339 12.125 9.13672C11.5417 9.72005 10.8333 10.0117 10 10.0117H6ZM6 8.01172H10C10.2833 8.01172 10.5208 7.91172 10.7125 7.71172C10.9042 7.51172 11 7.27839 11 7.01172C11 6.89505 10.9792 6.78672 10.9375 6.68672C10.8958 6.58672 10.8333 6.48672 10.75 6.38672C10.5667 6.17005 10.4458 5.95755 10.3875 5.74922C10.3292 5.54089 10.3 5.34505 10.3 5.16172C10.3 4.89505 10.3333 4.66589 10.4 4.47422C10.4667 4.28255 10.5 4.12839 10.5 4.01172C10.5 3.81172 10.4417 3.62839 10.325 3.46172C10.2083 3.29505 10.0583 3.17005 9.875 3.08672C9.725 3.02005 9.5875 2.94505 9.4625 2.86172C9.3375 2.77839 9.225 2.67005 9.125 2.53672C9.04167 2.43672 8.90417 2.32422 8.7125 2.19922C8.52083 2.07422 8.28333 2.01172 8 2.01172C7.71667 2.01172 7.47917 2.07839 7.2875 2.21172C7.09583 2.34505 6.95833 2.46172 6.875 2.56172C6.775 2.67839 6.6625 2.77839 6.5375 2.86172C6.4125 2.94505 6.275 3.02005 6.125 3.08672C5.94167 3.17005 5.79167 3.29505 5.675 3.46172C5.55833 3.62839 5.5 3.81172 5.5 4.01172C5.5 4.12839 5.53333 4.28255 5.6 4.47422C5.66667 4.66589 5.7 4.89505 5.7 5.16172C5.7 5.34505 5.67083 5.54089 5.6125 5.74922C5.55417 5.95755 5.43333 6.17005 5.25 6.38672C5.16667 6.48672 5.10417 6.58672 5.0625 6.68672C5.02083 6.78672 5 6.89505 5 7.01172C5 7.27839 5.09583 7.51172 5.2875 7.71172C5.47917 7.91172 5.71667 8.01172 6 8.01172ZM0 18.0117V15.2117C0 14.6451 0.145833 14.1242 0.4375 13.6492C0.729167 13.1742 1.11667 12.8117 1.6 12.5617C2.63333 12.0451 3.68333 11.6576 4.75 11.3992C5.81667 11.1409 6.9 11.0117 8 11.0117C9.1 11.0117 10.1833 11.1409 11.25 11.3992C12.3167 11.6576 13.3667 12.0451 14.4 12.5617C14.8833 12.8117 15.2708 13.1742 15.5625 13.6492C15.8542 14.1242 16 14.6451 16 15.2117V18.0117H0ZM2 16.0117H14V15.2117C14 15.0284 13.9542 14.8617 13.8625 14.7117C13.7708 14.5617 13.65 14.4451 13.5 14.3617C12.6 13.9117 11.6917 13.5742 10.775 13.3492C9.85833 13.1242 8.93333 13.0117 8 13.0117C7.06667 13.0117 6.14167 13.1242 5.225 13.3492C4.30833 13.5742 3.4 13.9117 2.5 14.3617C2.35 14.4451 2.22917 14.5617 2.1375 14.7117C2.04583 14.8617 2 15.0284 2 15.2117V16.0117Z\" fill=\"#00402A\"/>\n</svg>\n"
```

### `GET /assets/images/icons/social_distance_branco.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<path d=\"M6 21L2 17L6 13L7.4 14.4L5.8 16H18.2L16.6 14.4L18 13L22 17L18 21L16.6 19.6L18.2 18H5.8L7.4 19.6L6 21ZM2 11V10.425C2 10.025 2.10833 9.65833 2.325 9.325C2.54167 8.99167 2.84167 8.74167 3.225 8.575C3.65833 8.39167 4.10417 8.25 4.5625 8.15C5.02083 8.05 5.5 8 6 8C6.5 8 6.97917 8.05 7.4375 8.15C7.89583 8.25 8.34167 8.39167 8.775 8.575C9.15833 8.74167 9.45833 8.99167 9.675 9.325C9.89167 9.65833 10 10.025 10 10.425V11H2ZM14 11V10.425C14 10.025 14.1083 9.65833 14.325 9.325C14.5417 8.99167 14.8417 8.74167 15.225 8.575C15.6583 8.39167 16.1042 8.25 16.5625 8.15C17.0208 8.05 17.5 8 18 8C18.5 8 18.9792 8.05 19.4375 8.15C19.8958 8.25 20.3417 8.39167 20.775 8.575C21.1583 8.74167 21.4583 8.99167 21.675 9.325C21.8917 9.65833 22 10.025 22 10.425V11H14ZM6 7C5.45 7 4.97917 6.80417 4.5875 6.4125C4.19583 6.02083 4 5.55 4 5C4 4.45 4.19583 3.97917 4.5875 3.5875C4.97917 3.19583 5.45 3 6 3C6.55 3 7.02083 3.19583 7.4125 3.5875C7.80417 3.97917 8 4.45 8 5C8 5.55 7.80417 6.02083 7.4125 6.4125C7.02083 6.80417 6.55 7 6 7ZM18 7C17.45 7 16.9792 6.80417 16.5875 6.4125C16.1958 6.02083 16 5.55 16 5C16 4.45 16.1958 3.97917 16.5875 3.5875C16.9792 3.19583 17.45 3 18 3C18.55 3 19.0208 3.19583 19.4125 3.5875C19.8042 3.97917 20 4.45 20 5C20 5.55 19.8042 6.02083 19.4125 6.4125C19.0208 6.80417 18.55 7 18 7Z\" fill=\"#00A443\"/>\n</svg>\n"
```

### `GET /assets/images/icons/tiktok_silver.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"20\" height=\"21\" viewBox=\"0 0 20 21\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<circle cx=\"10\" cy=\"10.21\" r=\"10\" fill=\"#007F33\"/>\n<g clip-path=\"url(#clip0_2163_2040)\">\n<path d=\"M10.2477 4.60957C10.8601 4.60022 11.4679 4.60489 12.0756 4.60022C12.113 5.31546 12.3701 6.04473 12.8937 6.54961C13.4173 7.06851 14.1559 7.30693 14.8758 7.3864V9.27034C14.2026 9.24697 13.5248 9.10672 12.9124 8.81689C12.6459 8.69534 12.3981 8.54107 12.1551 8.38213C12.1504 9.74717 12.1597 11.1122 12.1457 12.4726C12.1083 13.127 11.8933 13.7768 11.5146 14.3144C10.9022 15.212 9.84103 15.7964 8.75181 15.8151C8.08331 15.8525 7.41481 15.6701 6.84449 15.3336C5.90018 14.7773 5.23636 13.7581 5.13819 12.6642C5.12884 12.4305 5.12416 12.1968 5.13351 11.9677C5.21766 11.0795 5.65709 10.2287 6.33961 9.649C7.11563 8.97583 8.20018 8.65327 9.21461 8.84493C9.22396 9.53681 9.19591 10.2287 9.19591 10.9205C8.73311 10.771 8.19083 10.813 7.78412 11.0935C7.48961 11.2852 7.26522 11.5797 7.14835 11.9116C7.05018 12.15 7.07823 12.4118 7.0829 12.6642C7.1951 13.4309 7.93372 14.076 8.71908 14.0059C9.24266 14.0012 9.74286 13.6974 10.014 13.2533C10.1028 13.099 10.201 12.9401 10.2057 12.7577C10.2524 11.921 10.2337 11.0888 10.2384 10.252C10.2431 8.36811 10.2337 6.48884 10.2477 4.60957Z\" fill=\"white\"/>\n</g>\n<defs>\n<clipPath id=\"clip0_2163_2040\">\n<rect width=\"11.2195\" height=\"11.2195\" fill=\"white\" transform=\"translate(4.39014 4.60022)\"/>\n</clipPath>\n</defs>\n</svg>\n"
```

### `GET /assets/images/icons/x-twitter_silver.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"20\" height=\"21\" viewBox=\"0 0 20 21\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<circle cx=\"10\" cy=\"10.21\" r=\"10\" fill=\"#007F33\"/>\n<g clip-path=\"url(#clip0_2163_2036)\">\n<path d=\"M10.9793 9.42295L14.8084 5.06769H13.9011L10.5748 8.84851L7.92007 5.06769H4.85742L8.87283 10.7855L4.85742 15.3522H5.76467L9.27515 11.3586L12.0793 15.3522H15.142L10.9793 9.42295ZM9.73628 10.8357L9.32882 10.266L6.09184 5.73694H7.48556L10.0987 9.39335L10.5044 9.96307L13.9007 14.7156H12.507L9.73628 10.8357Z\" fill=\"white\"/>\n</g>\n<defs>\n<clipPath id=\"clip0_2163_2036\">\n<rect width=\"11.2195\" height=\"11.2195\" fill=\"white\" transform=\"translate(4.39014 4.60022)\"/>\n</clipPath>\n</defs>\n</svg>\n"
```

### `GET /assets/images/icons/youtube_silver_v2.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"20\" height=\"21\" viewBox=\"0 0 20 21\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<g clip-path=\"url(#clip0_2163_2023)\">\n<g clip-path=\"url(#clip1_2163_2023)\">\n<path d=\"M10 20.21C7.69141 20.2116 5.45335 19.4143 3.66566 17.9535C1.87796 16.4927 0.650769 14.4584 0.192328 12.1958C-0.266113 9.93312 0.0724422 7.58155 1.15054 5.54011C2.22864 3.49868 3.97986 1.89315 6.10703 0.995987C8.55073 -0.0363683 11.3045 -0.0556829 13.7624 0.942292C16.2203 1.94027 18.1812 3.87378 19.2135 6.31749C20.2459 8.76119 20.2652 11.5149 19.2672 13.9729C18.2692 16.4308 16.3357 18.3916 13.892 19.424C12.6608 19.9451 11.337 20.2125 10 20.21ZM10 6.54299C8.60747 6.52613 7.21535 6.60263 5.83303 6.77199C5.60484 6.83661 5.39762 6.96005 5.23213 7.12993C5.06664 7.29981 4.94867 7.51018 4.89003 7.73999C4.73675 8.61018 4.66212 9.49241 4.66703 10.376C4.66191 11.2589 4.73621 12.1404 4.88903 13.01C4.94767 13.2398 5.06564 13.4502 5.23113 13.62C5.39662 13.7899 5.60384 13.9134 5.83203 13.978C7.21462 14.1484 8.60707 14.2259 10 14.21C11.3926 14.2269 12.7847 14.1503 14.167 13.981C14.3952 13.9164 14.6024 13.7929 14.7679 13.623C14.9334 13.4532 15.0514 13.2428 15.11 13.013C15.2634 12.1428 15.3381 11.2606 15.333 10.377C15.3381 9.4934 15.2634 8.61117 15.11 7.74099C15.0514 7.51118 14.9334 7.30081 14.7679 7.13093C14.6024 6.96105 14.3952 6.83762 14.167 6.77299C12.7847 6.60329 11.3926 6.52645 10 6.54299Z\" fill=\"#007F33\"/>\n<path d=\"M8.98389 12.404V9.01599L11.6939 10.71L8.98389 12.404Z\" fill=\"#007F33\"/>\n</g>\n</g>\n<defs>\n<clipPath id=\"clip0_2163_2023\">\n<rect width=\"20\" height=\"20\" fill=\"white\" transform=\"translate(0 0.209961)\"/>\n</clipPath>\n<clipPath id=\"clip1_2163_2023\">\n<rect width=\"20\" height=\"20\" fill=\"white\" transform=\"translate(0 0.209961)\"/>\n</clipPath>\n</defs>\n</svg>\n"
```

### `GET /assets/images/logo-aneel.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"439\" height=\"47\" viewBox=\"0 0 439 47\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M98.8792 39.2805L97.4819 46.9153H101.352L101.481 46.2174H98.4919L99.0297 43.2779H101.136L101.263 42.5799H99.1583L99.6346 39.9784H102.486L102.614 39.2805H98.8792Z\" fill=\"#005172\"/>\n<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M104.445 40.8072L103.327 46.9148H104.034L105.007 41.6023H105.024L106.079 46.9148H106.84L107.958 40.8072H107.252L106.428 45.3068H106.411L105.55 40.8072H104.445Z\" fill=\"#005172\"/>\n<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M109.912 40.8072L108.794 46.9148H111.891L111.994 46.3565H109.603L110.033 44.004H111.718L111.82 43.4457H110.136L110.516 41.3656H112.799L112.9 40.8072H109.912Z\" fill=\"#005172\"/>\n<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M115.386 41.3658H116.002C117.016 41.3658 117.273 41.7381 117.147 42.4223C117.023 43.1079 116.628 43.4801 115.614 43.4801H114.999L115.386 41.3658ZM113.664 46.915H114.37L114.893 44.0549L115.822 44.0221L116.344 46.915H117.087L116.535 43.8866C117.303 43.6498 117.736 43.0751 117.86 42.3895C118.045 41.3823 117.407 40.8075 116.357 40.8075H114.782L113.664 46.915Z\" fill=\"#005172\"/>\n<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M123.238 42.6852C123.327 41.704 123.072 40.7228 121.868 40.7228C120.718 40.7228 119.841 41.5521 119.418 43.8607C118.997 46.1707 119.569 47 120.718 47C121.252 47 121.671 46.7961 122.005 46.4581L122.012 46.9152H122.463L122.987 44.055H121.393L121.292 44.6147H122.179L122.021 45.4769C121.906 46.1023 121.482 46.4417 120.821 46.4417C120.033 46.4417 119.774 45.7807 120.126 43.8607C120.476 41.9407 120.978 41.2811 121.765 41.2811C122.481 41.2811 122.623 41.8901 122.532 42.6852H123.238Z\" fill=\"#005172\"/>\n<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M125.591 40.8072L124.473 46.9148H125.179L126.297 40.8072H125.591Z\" fill=\"#005172\"/>\n<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M129.715 41.2643H129.732L129.857 44.5459H128.39L129.715 41.2643ZM129.318 40.8072L126.761 46.9148H127.466L128.187 45.1043H129.854L129.912 46.9148H130.618L130.296 40.8072H129.318Z\" fill=\"#005172\"/>\n<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M136.548 39.2805L135.151 46.9153H139.021L139.15 46.2174H136.161L136.7 43.2779H138.805L138.932 42.5799H136.828L137.304 39.9784H140.156L140.283 39.2805H136.548Z\" fill=\"#005172\"/>\n<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M142.206 40.8072L141.088 46.9148H143.912L144.015 46.3565H141.896L142.912 40.8072H142.206Z\" fill=\"#005172\"/>\n<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M149.393 39.3183H148.687L147.587 40.4857H148.039L149.393 39.3183ZM146.577 40.8072L145.459 46.9148H148.556L148.657 46.3565H146.266L146.697 44.004H148.382L148.483 43.4471H146.8L147.181 41.3656H149.462L149.564 40.8072H146.577Z\" fill=\"#005172\"/>\n<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M6.9224 39.8325H6.94566L7.09893 43.9352H5.26653L6.9224 39.8325ZM6.42701 39.2619L3.23022 46.8953H4.1129L5.01336 44.6332H7.09619L7.16872 46.8953H8.0514L7.64906 39.2619H6.42701Z\" fill=\"#005172\"/>\n<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M13.9964 42.6663C14.0853 41.6851 13.8308 40.7039 12.6265 40.7039C11.4756 40.7039 10.5998 41.5332 10.1783 43.8418C9.75542 46.1518 10.3274 46.9797 11.477 46.9797C12.0121 46.9797 12.4294 46.7772 12.762 46.4392L12.7702 46.8963H13.2218L13.7459 44.0361H12.1516L12.0504 44.5945H12.9371L12.7798 45.458C12.6648 46.0834 12.2406 46.4214 11.5796 46.4214C10.7914 46.4214 10.5327 45.7618 10.8844 43.8418C11.2348 41.9218 11.737 41.2622 12.5239 41.2622C13.2396 41.2622 13.3819 41.8712 13.2902 42.6663H13.9964Z\" fill=\"#005172\"/>\n<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M16.5875 40.4669H17.1664L17.9683 39.6047L18.4528 40.4669H19.0316L18.331 39.2995H17.7151L16.5875 40.4669ZM16.2837 40.7884L15.1656 46.896H18.2625L18.3652 46.3377H15.9744L16.4055 43.9866H18.0887L18.1914 43.4269H16.5068L16.8886 41.3468H19.1698L19.2725 40.7884H16.2837Z\" fill=\"#005172\"/>\n<path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M21.0494 40."
```

### `GET /assets/images/logo-conexao-digital.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"198\" height=\"68\" viewBox=\"0 0 198 68\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<path d=\"M160.344 15.5331C157.139 15.5331 154.767 13.2042 154.767 9.82836C154.767 6.47391 157.139 4.14502 160.344 4.14502C163.548 4.14502 165.92 6.47391 165.92 9.82836C165.92 13.2042 163.548 15.5331 160.344 15.5331ZM160.344 13.4178C162.16 13.4178 163.506 12.0932 163.506 9.82836C163.506 7.58494 162.16 6.28161 160.344 6.28161C158.527 6.28161 157.203 7.58494 157.203 9.82836C157.203 12.0932 158.527 13.4178 160.344 13.4178Z\" fill=\"#00A552\"/>\n<path d=\"M147.391 1.83747C147.006 1.83747 146.75 2.17933 146.75 2.8203H144.827C144.827 1.02556 146.024 0 147.241 0C148.716 0 149.079 1.0683 149.763 1.0683C150.147 1.0683 150.382 0.726442 150.382 0.10683H152.326C152.326 1.90157 151.109 2.90577 149.912 2.90577C148.438 2.90577 148.075 1.83747 147.391 1.83747ZM148.31 4.14499C150.831 4.14499 153.117 5.46968 153.117 8.46091V15.298H150.895V14.2297C150.233 15.0844 149.079 15.533 147.583 15.533C145.126 15.533 143.31 14.2938 143.31 12.1572C143.31 10.042 145.169 8.76004 147.733 8.76004C149.079 8.76004 150.126 9.14462 150.767 9.82833V8.69594C150.767 6.98666 149.698 6.26022 148.224 6.26022C146.9 6.26022 145.981 6.98666 145.596 8.16179L143.609 7.37125C144.272 5.14919 146.13 4.14499 148.31 4.14499ZM148.288 13.7597C149.763 13.7597 150.81 13.1828 150.81 12.1359C150.81 11.0676 149.763 10.4907 148.288 10.4907C146.835 10.4907 145.746 11.0462 145.746 12.1359C145.746 13.2042 146.835 13.7597 148.288 13.7597Z\" fill=\"#00A552\"/>\n<path d=\"M140.42 4.38H143.176L139.159 9.72149L143.326 15.298H140.505L137.749 11.6017L134.972 15.298H132.215L136.382 9.76422L132.365 4.38H135.185L137.792 7.88402L140.42 4.38Z\" fill=\"#00A552\"/>\n<path d=\"M132.351 9.76426C132.351 9.95655 132.351 10.1702 132.33 10.3625H124.382C124.553 12.4777 125.706 13.4819 127.288 13.4819C128.655 13.4819 129.702 12.8837 130.321 11.2812L132.33 12.0504C131.411 14.5075 129.488 15.5331 127.33 15.5331C124.211 15.5331 121.989 13.2469 121.989 9.82836C121.989 6.40981 124.189 4.14502 127.266 4.14502C130.236 4.14502 132.351 6.40981 132.351 9.76426ZM127.245 6.11069C125.771 6.11069 124.766 6.98669 124.467 8.6746H129.916C129.616 6.92259 128.612 6.11069 127.245 6.11069Z\" fill=\"#00A552\"/>\n<path d=\"M116.391 4.14502C118.677 4.14502 120.344 5.81156 120.344 8.35411V15.298H117.972V8.7387C117.972 7.17898 117.053 6.23888 115.771 6.23888C114.297 6.23888 113.25 7.22171 113.25 8.84553V15.298H110.878V4.38004H113.25V5.91839C113.934 4.6578 115.109 4.14502 116.391 4.14502Z\" fill=\"#00A552\"/>\n<path d=\"M103.57 15.5331C100.365 15.5331 97.9932 13.2042 97.9932 9.82836C97.9932 6.47391 100.365 4.14502 103.57 4.14502C106.775 4.14502 109.146 6.47391 109.146 9.82836C109.146 13.2042 106.775 15.5331 103.57 15.5331ZM103.57 13.4178C105.386 13.4178 106.732 12.0932 106.732 9.82836C106.732 7.58494 105.386 6.28161 103.57 6.28161C101.754 6.28161 100.429 7.58494 100.429 9.82836C100.429 12.0932 101.754 13.4178 103.57 13.4178Z\" fill=\"#00A552\"/>\n<path d=\"M90.4284 15.533C86.0697 15.533 82.9717 12.4136 82.9717 7.81992C82.9717 3.22624 86.0697 0.106812 90.4284 0.106812C93.6119 0.106812 96.069 1.79472 97.0732 4.76459L94.7657 5.64059C94.0606 3.63219 92.4795 2.41433 90.4284 2.41433C87.6295 2.41433 85.5783 4.4441 85.5783 7.81992C85.5783 11.1957 87.6295 13.2255 90.4284 13.2255C92.4795 13.2255 94.0606 12.0076 94.7657 9.99924L97.0732 10.8752C96.069 13.8451 93.6119 15.533 90.4284 15.533Z\" fill=\"#00A552\"/>\n<path d=\"M39.0265 20.2067L42.1577 23.3379V42.0342L39.0265 45.1654H11.0273V32.686L17.2443 26.469V38.903H39.0265L35.9407 35.8172V26.469H11.0273V20.2067H39.0265Z\" fill=\"#F59A26\"/>\n<path d=\"M51.5058 20.2974H45.2888V45.2107H51.5058V20.2974Z\" fill=\"#F59A26\"/>\n<path d=\"M54.637 42.0795V20.2974H82.6361L85.7673 23.3832V26.5144H60.854V38.9937H79.5503V35.8625H63.9398V29.6455H85.7673V45.2107H57.7228L54.637 42.0795Z\" fill=\"#F59A26\"/>\n<path d=\"M95.1154 20.2974H88.8984V45.2107H95.1154V20.2974Z\" fill=\"#F59A26\"/>\n<path d=\"M110.681 26.5144H98.2012V20.2974H129.377V26.5144H116.898V45.2107H110.681V26.51"
```

### `GET /assets/images/logo_neoenergia_letra_branca.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"150\" height=\"40\" viewBox=\"0 0 150 40\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<path d=\"M51.2022 21.8995V23.8826C50.7178 23.2068 50.3303 22.671 50.0397 22.275L42.6608 12.2969H40.6045V27.715H42.8723V16.1293C43.3575 16.8051 43.7449 17.3409 44.0348 17.7369L51.4137 27.715H53.4619V12.2969H51.2022V21.8995Z\" fill=\"white\"/>\n<path d=\"M62.8924 17.0215C62.0893 16.5371 61.1659 16.2889 60.2282 16.3055C59.4849 16.2961 58.7473 16.4343 58.0579 16.7123C57.4115 16.9805 56.8298 17.3835 56.3515 17.8945C55.8573 18.4237 55.4784 19.0499 55.2391 19.7334C54.9728 20.5016 54.8422 21.3103 54.8532 22.1233C54.8532 23.2857 55.0884 24.3024 55.5588 25.1735C56.0036 26.0217 56.6742 26.7305 57.4965 27.2216C58.3579 27.7246 59.3413 27.98 60.3386 27.9597C60.976 27.9609 61.6101 27.8681 62.2205 27.6843C62.8337 27.4993 63.4031 27.192 63.8944 26.7811C64.4289 26.3286 64.8737 25.7797 65.2056 25.163L63.5317 24.1273C63.3239 24.5466 63.0401 24.9238 62.6948 25.2397C62.3824 25.5194 62.0144 25.7298 61.6149 25.857C61.2008 25.9844 60.7695 26.0475 60.3363 26.0441C59.7415 26.0529 59.1555 25.8994 58.6415 25.6001C58.1272 25.2894 57.7173 24.8325 57.4639 24.2877C57.2208 23.765 57.0815 23.2001 57.0536 22.6243H65.3834C65.3974 22.5218 65.4044 22.4185 65.4043 22.3151V22.0292C65.4043 20.8838 65.1804 19.8818 64.7325 19.0232C64.3159 18.1987 63.6789 17.5059 62.8924 17.0215ZM57.8824 19.158C58.1449 18.8195 58.4899 18.5542 58.8844 18.3873C59.3028 18.2139 59.752 18.1273 60.2049 18.1328C60.6161 18.1222 61.0244 18.2038 61.3999 18.3715C61.7754 18.5393 62.1086 18.7889 62.3751 19.1022C62.7789 19.5672 63.0385 20.2061 63.1539 21.019H57.1176C57.1501 20.8144 57.1959 20.6121 57.2547 20.4134C57.3819 19.9587 57.595 19.5326 57.8824 19.158Z\" fill=\"white\"/>\n<path d=\"M74.9278 17.0323C74.0327 16.5366 73.0227 16.286 71.9997 16.3058C71.2307 16.2985 70.4676 16.4406 69.7528 16.7242C69.0743 16.9949 68.4595 17.4039 67.9476 17.925C67.4282 18.4533 67.023 19.0828 66.7573 19.7744C66.4718 20.5231 66.3298 21.3189 66.3388 22.1201C66.3388 23.2655 66.5848 24.2787 67.0769 25.1598C67.5486 26.0204 68.25 26.7331 69.103 27.2184C69.889 27.644 70.7587 27.8922 71.651 27.9456C72.5432 27.9991 73.4364 27.8564 74.2676 27.5276C74.9442 27.2508 75.5593 26.8427 76.0774 26.3269C76.5976 25.8034 77.0026 25.1771 77.2666 24.488C77.5535 23.732 77.6954 22.9287 77.685 22.1201C77.685 20.9577 77.439 19.941 76.9469 19.07C76.4754 18.217 75.7764 17.5116 74.9278 17.0323ZM75.2149 23.7498C75.0581 24.2027 74.8105 24.6188 74.4873 24.9727C74.1797 25.3025 73.8041 25.5615 73.3864 25.7317C72.9464 25.9112 72.475 26.0013 71.9997 25.9968C71.3843 26.0042 70.7786 25.8435 70.248 25.5318C69.7182 25.214 69.2903 24.7514 69.0146 24.1985C68.714 23.6181 68.5636 22.9245 68.5636 22.1178C68.5535 21.5648 68.6392 21.0142 68.817 20.4904C68.9724 20.0429 69.2159 19.6311 69.5331 19.2792C69.8362 18.944 70.2131 18.6838 70.6339 18.519C71.0689 18.3485 71.5325 18.2625 71.9997 18.2656C72.6291 18.2546 73.2496 18.4153 73.7945 18.7306C74.3231 19.0447 74.7519 19.5021 75.0313 20.0499C75.3323 20.6218 75.4834 21.3134 75.4834 22.1201C75.4902 22.6749 75.3993 23.2265 75.2149 23.7498Z\" fill=\"white\"/>\n<path d=\"M86.6554 17.0215C85.8523 16.5371 84.9289 16.2889 83.9911 16.3055C83.248 16.2962 82.5103 16.4344 81.8209 16.7123C81.1749 16.9806 80.5936 17.3837 80.1157 17.8945C79.621 18.4234 79.2422 19.0497 79.0033 19.7334C78.7369 20.5016 78.6063 21.3103 78.6173 22.1233C78.6173 23.2857 78.8525 24.3024 79.3229 25.1735C79.7678 26.0217 80.4384 26.7305 81.2607 27.2216C82.121 27.7227 83.1027 27.9768 84.0981 27.9563C84.7363 27.9574 85.3712 27.8646 85.9824 27.6808C86.5952 27.495 87.1645 27.1879 87.6562 26.7776C88.1901 26.3247 88.6344 25.7759 88.9662 25.1595L87.2924 24.1238C87.0841 24.5429 86.8004 24.92 86.4555 25.2362C86.1429 25.5157 85.7749 25.726 85.3756 25.8535C84.9614 25.9808 84.5302 26.0439 84.0969 26.0406C83.5021 26.0503 82.9159 25.898 82.401 25.6001C81.8885 25.2887 81.4802 24.8318 81.2281 24.2877C80.9845 23.7652 80.8451 23.2002 80.8178 22.6243H89.1453C89.1602 22.5219 89.1675 22.4186 89.1673 22."
```

### `GET /assets/images/logo_v2.svg`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<svg width=\"116\" height=\"25\" viewBox=\"0 0 116 25\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<path d=\"M39.2295 14.019V15.6158C38.8395 15.0717 38.5275 14.6402 38.2935 14.3214L32.352 6.28705H30.6963V18.7017H32.5224V9.37292C32.913 9.91703 33.225 10.3485 33.4583 10.6674L39.3998 18.7017H41.049V6.28705H39.2295V14.019Z\" fill=\"#007F33\"/>\n<path d=\"M48.6426 10.0918C47.9959 9.70178 47.2524 9.50195 46.4974 9.51527C45.8989 9.5077 45.305 9.61905 44.7499 9.84286C44.2294 10.0588 43.761 10.3833 43.3759 10.7947C42.9779 11.2209 42.6729 11.7251 42.4802 12.2754C42.2658 12.894 42.1606 13.5452 42.1695 14.1998C42.1695 15.1357 42.3588 15.9544 42.7376 16.6557C43.0958 17.3387 43.6357 17.9095 44.2978 18.3049C44.9914 18.7099 45.7833 18.9156 46.5863 18.8992C47.0995 18.9002 47.6101 18.8254 48.1016 18.6774C48.5953 18.5285 49.0538 18.2811 49.4494 17.9502C49.8798 17.5859 50.238 17.1439 50.5052 16.6473L49.1574 15.8134C48.99 16.151 48.7615 16.4547 48.4835 16.7091C48.232 16.9343 47.9356 17.1037 47.614 17.2061C47.2805 17.3087 46.9333 17.3595 46.5844 17.3568C46.1055 17.3639 45.6337 17.2403 45.2198 16.9992C44.8057 16.7491 44.4756 16.3812 44.2716 15.9425C44.0759 15.5217 43.9637 15.0668 43.9412 14.6032H50.6484C50.6596 14.5207 50.6652 14.4375 50.6652 14.3542V14.124C50.6652 13.2017 50.4849 12.3949 50.1242 11.7036C49.7888 11.0397 49.2759 10.4818 48.6426 10.0918ZM44.6086 11.8121C44.8199 11.5396 45.0977 11.3259 45.4154 11.1916C45.7523 11.0519 46.114 10.9822 46.4786 10.9866C46.8097 10.9781 47.1385 11.0438 47.4408 11.1788C47.7432 11.3139 48.0115 11.5149 48.2261 11.7672C48.5512 12.1416 48.7602 12.6561 48.8532 13.3106H43.9927C44.0189 13.1458 44.0558 12.9829 44.1032 12.823C44.2056 12.4569 44.3772 12.1137 44.6086 11.8121Z\" fill=\"#007F33\"/>\n<path d=\"M58.3336 10.1002C57.6129 9.70108 56.7996 9.49929 55.9759 9.51522C55.3567 9.50936 54.7423 9.62379 54.1667 9.85216C53.6204 10.0701 53.1254 10.3994 52.7132 10.819C52.295 11.2444 51.9687 11.7513 51.7548 12.3081C51.5249 12.911 51.4106 13.5517 51.4178 14.1969C51.4178 15.1191 51.6159 15.935 52.0121 16.6445C52.3919 17.3374 52.9567 17.9113 53.6435 18.302C54.2764 18.6447 54.9767 18.8446 55.6952 18.8876C56.4136 18.9306 57.1327 18.8157 57.802 18.551C58.3468 18.3281 58.8421 17.9995 59.2593 17.5842C59.6781 17.1627 60.0043 16.6583 60.2168 16.1035C60.4478 15.4947 60.5621 14.8479 60.5537 14.1969C60.5537 13.2609 60.3556 12.4423 59.9594 11.7409C59.5798 11.0541 59.017 10.4861 58.3336 10.1002ZM58.5648 15.5091C58.4385 15.8738 58.2392 16.2088 57.9789 16.4938C57.7313 16.7594 57.4288 16.9679 57.0926 17.1049C56.7382 17.2495 56.3586 17.322 55.9759 17.3183C55.4805 17.3244 54.9927 17.1949 54.5654 16.944C54.1389 16.688 53.7944 16.3156 53.5724 15.8704C53.3303 15.4031 53.2092 14.8446 53.2092 14.195C53.2011 13.7497 53.2701 13.3064 53.4133 12.8847C53.5384 12.5243 53.7344 12.1927 53.9898 11.9094C54.2339 11.6395 54.5373 11.43 54.8762 11.2973C55.2265 11.16 55.5997 11.0907 55.9759 11.0933C56.4827 11.0844 56.9823 11.2138 57.4211 11.4676C57.8467 11.7206 58.192 12.0889 58.4169 12.53C58.6594 12.9905 58.781 13.5473 58.781 14.1969C58.7865 14.6436 58.7133 15.0878 58.5648 15.5091Z\" fill=\"#007F33\"/>\n<path d=\"M67.7764 10.0918C67.1298 9.70178 66.3862 9.50195 65.6312 9.51527C65.0328 9.50778 64.4388 9.61912 63.8837 9.84286C63.3636 10.0589 62.8955 10.3834 62.5107 10.7947C62.1124 11.2206 61.8073 11.725 61.615 12.2754C61.4005 12.894 61.2954 13.5452 61.3042 14.1998C61.3042 15.1357 61.4936 15.9544 61.8724 16.6557C62.2305 17.3387 62.7705 17.9095 63.4326 18.3049C64.1253 18.7083 64.9158 18.913 65.7173 18.8964C66.2312 18.8974 66.7424 18.8226 67.2345 18.6746C67.728 18.5251 68.1863 18.2777 68.5823 17.9474C69.0121 17.5827 69.3699 17.1408 69.6371 16.6445L68.2893 15.8106C68.1217 16.148 67.8932 16.4517 67.6154 16.7063C67.3638 16.9313 67.0675 17.1006 66.7459 17.2033C66.4124 17.3058 66.0652 17.3566 65.7164 17.354C65.2374 17.3618 64.7654 17.2392 64.3508 16.9992C63.9381 16.7485 63.6094 16.3807 63.4064 15.9425C63.2102 15.5218 63.098 15.0669 63.076 14.6032H69.7812C69.7933 14.5207 69.7992 14.4375 69.7"
```

### `GET /bandeira-tarifaria/1.0.0/bandeira-tarifaria`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "canalSolicitante": "AGC",
  "usuario": "***",
  "distribuidora": "COELBA"
}
```
response body:
```json
{
  "bandeiraTarifaria": [
    {
      "bandeira": "VD",
      "descricao": "Bandeira Verde",
      "mesReferencia": "01/2026",
      "dataInclusao": "01/01/2026",
      "validoAte": "31/12/9999"
    }
  ],
  "e_resultado": "X",
  "retorno": {
    "tipo": "S",
    "id": "ZATCWS",
    "numero": "029",
    "mensagem": "Executado com sucesso."
  }
}
```

### `GET /imoveis/1.1.0/clientes/03021937586/ucs`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "documento": "*******7586",
  "canalSolicitante": "AGC",
  "distribuidora": "COELBA",
  "usuario": "***",
  "indMaisUcs": "X",
  "protocolo": "123",
  "opcaoSSOS": "S",
  "tipoPerfil": "1"
}
```
response body:
```json
{
  "ucs": [
    {
      "status": "LIGADA",
      "uc": "007098175908",
      "nomeCliente": "PAULA PEREIRA FERNANDES",
      "instalacao": "0009043917",
      "local": {
        "endereco": "RUA BAHIA, 130",
        "bairro": "CENTRO-LAPAO",
        "municipio": "LAPAO",
        "cep": "44905-000",
        "uf": "BA"
      },
      "enderecoEntrega": {
        "endereco": "RUA BAHIA 130",
        "bairro": "CENTRO-LAPAO",
        "municipio": "LAPAO",
        "cep": "44905-000",
        "uf": "BA"
      },
      "isGrupo": null,
      "nomeGrupo": null,
      "grupoTensao": "B",
      "bOptante": null,
      "contrato": "0023513361",
      "dt_inicio": "2026-02-20",
      "dt_fim": "9999-12-31",
      "medidor": null,
      "indMaisUcs": null,
      "ucColetiva": null,
      "indCCColetiva": null
    },
    {
      "status": "DESLIGADA",
      "uc": "007085489032",
      "nomeCliente": "PAULA PEREIRA FERNANDES",
      "instalacao": "0009043917",
      "local": {
        "endereco": "RUA BAHIA, 130",
        "bairro": "CENTRO-LAPAO",
        "municipio": "LAPAO",
        "cep": "44905-000",
        "uf": "BA"
      },
      "enderecoEntrega": {
        "endereco": "RUA BAHIA 130",
        "bairro": "CENTRO-LAPAO",
        "municipio": "LAPAO",
        "cep": "44905-000",
        "uf": "BA"
      },
      "isGrupo": null,
      "nomeGrupo": null,
      "grupoTensao": "B",
      "bOptante": null,
      "contrato": "0022215892",
      "dt_inicio": "2024-07-31",
      "dt_fim": "2024-12-19",
      "medidor": null,
      "indMaisUcs": null,
      "ucColetiva": null,
      "indCCColetiva": null
    }
  ],
  "e_resultado": "X",
  "retorno": {
    "id": "ZATCWS",
    "tipo": "S",
    "numero": "029",
    "mensagem": "Executado com sucesso."
  }
}
```

### `GET /main.12d1bbaf63a33bd2.js`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"(self[\"webpackChunkneoenergia_web\"] = self[\"webpackChunkneoenergia_web\"] || []).push([[179],{\n\n/***/ 20721:\n/*!**********************************!*\\\n  !*** ./src/app/app.component.ts ***!\n  \\**********************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n\"use strict\";\n__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"AppComponent\": () => (/* binding */ AppComponent)\n/* harmony export */ });\n/* harmony import */ var _angular_router__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(/*! @angular/router */ 74202);\n/* harmony import */ var _environments_environment__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @environments/environment */ 24766);\n/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! @angular/core */ 5000);\n/* harmony import */ var _angular_platform_browser__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(/*! @angular/platform-browser */ 22313);\n/* harmony import */ var _core_services_customsweetalert_loading_service__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./core/services/customsweetalert/loading.service */ 74139);\n/* harmony import */ var _angular_common__WEBPACK_IMPORTED_MODULE_8__ = __webpack_require__(/*! @angular/common */ 69808);\n/* harmony import */ var _shared_components_spinner_spinner_component__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./shared/components/spinner/spinner.component */ 77477);\n/* harmony import */ var _core_header_header_component__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./core/header/header.component */ 6885);\n/* harmony import */ var _core_footer_footer_component__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ./core/footer/footer.component */ 61594);\n\n\n\n\n\n\n\n\n\n\n\nfunction AppComponent_app_spinner_0_Template(rf, ctx) {\n  if (rf & 1) {\n    _angular_core__WEBPACK_IMPORTED_MODULE_5__[\"ɵɵelement\"](0, \"app-spinner\", 1);\n  }\n\n  if (rf & 2) {\n    const ctx_r0 = _angular_core__WEBPACK_IMPORTED_MODULE_5__[\"ɵɵnextContext\"]();\n    _angular_core__WEBPACK_IMPORTED_MODULE_5__[\"ɵɵproperty\"](\"mensagem\", ctx_r0.loading.textLabel);\n  }\n}\n\nlet AppComponent = /*#__PURE__*/(() => {\n  class AppComponent {\n    constructor(_router, activatedRoute, titleService, cd, loading) {\n      this._router = _router;\n      this.activatedRoute = activatedRoute;\n      this.titleService = titleService;\n      this.cd = cd;\n      this.loading = loading;\n      this.loadingDataImg = false;\n      this.title = 'Agência Virtual';\n      this.activatedRoute = this.activatedRoute;\n      this.titleService = this.titleService;\n      const script = document.createElement('script');\n      script.async = true;\n      script.src = 'https://www.googletagmanager.com/gtag/js?id=' + _environments_environment__WEBPACK_IMPORTED_MODULE_0__.environment.gtag;\n      document.head.prepend(script);\n      gtag('config', _environments_environment__WEBPACK_IMPORTED_MODULE_0__.environment.gtag); // const navigation = this._router.events.pipe(\n      //     filter((e): e is NavigationEnd => e instanceof NavigationEnd)\n      // );\n      // navigation.subscribe((e: NavigationEnd) => {\n      //     const lastSegment = e.urlAfterRedirects.split('/').pop();\n      //     gtag('config', environment.gtag, { 'page_path': lastSegment });\n      // });\n    }\n\n    ngAfterContentChecked() {\n      this.cd.detectChanges();\n    }\n\n    getTitle(state, parent) {\n      var data = [];\n\n      if (parent && parent.snapshot.data && parent.snapshot.data.title) {\n        data.push(parent.snapshot.data.title);\n      }\n\n      if (state && parent) {\n        data.push(...this.getTitle(state, state.firstChild(parent)));\n      }\n\n      return data;\n    }\n\n    navigationInterceptor(event) {\n      if (event instanceof _angular_router__WEBPACK_IMPORTED_MODULE_6__.NavigationStart) {\n        this.loadingDataImg = true;\n      }\n\n      if (event instanceof _angular_rou"
```

### `GET /material-icons-outlined.f86cb7b0aa53f0fe.woff2`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"<unreadable: 'utf-8' codec can't decode byte 0x8c in position 11: invalid start byte>"
```

### `GET /multilogin/2.0.0/agv/cliente/03021937586/COELBA/grupo-de-cliente`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "tipoPerfil": "0"
}
```
response body:
```json
{
  "dados": [
    {
      "funcionalidades": [
        {
          "codigo": "F035"
        },
        {
          "codigo": "F000"
        },
        {
          "codigo": "F041"
        },
        {
          "codigo": "F069"
        },
        {
          "codigo": "F023"
        },
        {
          "codigo": "F013"
        },
        {
          "codigo": "F088"
        },
        {
          "codigo": "F003"
        },
        {
          "codigo": "F066"
        },
        {
          "codigo": "F009"
        },
        {
          "codigo": "F060"
        },
        {
          "codigo": "F089"
        },
        {
          "codigo": "F021"
        },
        {
          "codigo": "F062"
        },
        {
          "codigo": "F087"
        },
        {
          "codigo": "F027"
        },
        {
          "codigo": "F005"
        },
        {
          "codigo": "F059"
        },
        {
          "codigo": "F065"
        },
        {
          "codigo": "F075"
        },
        {
          "codigo": "F067"
        },
        {
          "codigo": "F015"
        },
        {
          "codigo": "F011"
        },
        {
          "codigo": "F007"
        },
        {
          "codigo": "F063"
        },
        {
          "codigo": "F076"
        },
        {
          "codigo": "F057"
        },
        {
          "codigo": "F081"
        },
        {
          "codigo": "F086"
        },
        {
          "codigo": "F072"
        },
        {
          "codigo": "F022"
        },
        {
          "codigo": "F002"
        },
        {
          "codigo": "F080"
        },
        {
          "codigo": "F058"
        },
        {
          "codigo": "F056"
        },
        {
          "codigo": "F061"
        },
        {
          "codigo": "F085"
        },
        {
          "codigo": "F068"
        },
        {
          "codigo": "F004"
        },
        {
          "codigo": "F064"
        },
        {
          "codigo": "F012"
        },
        {
          "codigo": "F006"
        },
        {
          "codigo": "F008"
        },
        {
          "codigo": "F001"
        },
        {
          "codigo": "F055"
        }
      ],
      "codigo": 1,
      "nome": "Comum",
      "ativo": true
    }
  ],
  "retorno": {
    "tipo": "Backoffice MultiLogin",
    "id": 200,
    "numero": "200_OK",
    "mensagem": "Sucesso"
  }
}
```

### `GET /multilogin/2.0.0/servicos/datacerta/ucs/007085489032/datacerta`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "codigo": "007085489032",
  "canalSolicitante": "AGC",
  "usuario": "***",
  "operacao": "CON",
  "tipoPerfil": "1",
  "documentoSolicitante": "***",
  "distribuidora": "COELBA"
}
```
response body:
```json
{
  "e_resultado": "X",
  "dataAtual": "Data Normal Regulada",
  "retorno": {
    "id": "ZATCWS",
    "tipo": "S",
    "numero": "029",
    "mensagem": "Executado com sucesso."
  }
}
```

### `GET /multilogin/2.0.0/servicos/datacerta/ucs/007098175908/datacerta`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "codigo": "007098175908",
  "canalSolicitante": "AGC",
  "usuario": "***",
  "operacao": "CON",
  "tipoPerfil": "1",
  "documentoSolicitante": "***",
  "distribuidora": "COELBA"
}
```
response body:
```json
{
  "e_resultado": "X",
  "possuiDataBoa": "X",
  "dataAtual": "06",
  "retorno": {
    "id": "ZATCWS",
    "tipo": "S",
    "numero": "029",
    "mensagem": "Executado com sucesso."
  }
}
```

### `GET /multilogin/2.0.0/servicos/debito-automatico/conta-cadastrada-debito`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "codigo": "007098175908",
  "codCliente": "1014945628",
  "canalSolicitante": "AGC",
  "usuario": "***",
  "valida": "",
  "distribuidora": "COELBA",
  "tipoPerfil": "1",
  "documentoSolicitante": "***"
}
```
response body:
```json
{
  "e_resultado": "X",
  "retorno": {
    "id": "ZATCWS",
    "tipo": "S",
    "numero": "145",
    "mensagem": "UC não possui débito automático cadastrado"
  }
}
```

### `GET /multilogin/2.0.0/servicos/fatura-digital/ucs/fatura-digital`
- chamadas observadas: `2`
- status HTTP: `[200, 400]`
```json
{
  "codigo": "007098175908",
  "canalSolicitante": "AGC",
  "usuario": "***",
  "distribuidora": "COELBA",
  "tipoPerfil": "1"
}
```
response body:
```json
{
  "e_resultado": "X",
  "PossuiFaturaDigital": null,
  "emailFatura": null,
  "emailCadastro": "paulinha_fernandes7@hotmail.com",
  "dominioWhatsapp": null,
  "dominioSMS": null,
  "faturaBraile": null,
  "faturaEntregaAlternativa": null,
  "retorno": {
    "tipo": "S",
    "id": "ZATCWS",
    "numero": "029",
    "mensagem": "Executado com sucesso."
  }
}
```

### `GET /multilogin/2.0.0/servicos/faturas/334075735546/dados-pagamento`
- chamadas observadas: `1`
- status HTTP: `[400]`
```json
{
  "codigo": "007085489032",
  "protocolo": "20260416280227433",
  "usuario": "***",
  "canalSolicitante": "AGC",
  "distribuidora": "COELBA",
  "regiao": "NE",
  "tipoPerfil": "1",
  "byPassActiv": "X",
  "documentoSolicitante": "*******7586",
  "documento": "*******7586"
}
```
response body:
```json
{
  "retorno": {
    "tipo": "E",
    "id": "ZATCWS",
    "numero": "448",
    "mensagem": "Fatura vencida a mais de 365 dias. Favor procurar demais canais."
  }
}
```

### `GET /multilogin/2.0.0/servicos/faturas/334075735546/pdf`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "codigo": "007085489032",
  "protocolo": "20260416280227433",
  "tipificacao": "1031602",
  "usuario": "***",
  "canalSolicitante": "AGC",
  "motivo": "02",
  "distribuidora": "COELBA",
  "regiao": "NE",
  "tipoPerfil": "1",
  "documento": "*******7586",
  "documentoSolicitante": "*******7586",
  "byPassActiv": ""
}
```
response body:
```json
{
  "fileName": "007085489032",
  "fileSize": "51724 ",
  "fileData": "<base64:68968 chars>",
  "fileExtension": ".pdf",
  "retorno": {
    "tipo": "S",
    "id": "ZATCWS",
    "numero": "066",
    "mensagem": "Documento gerado com sucesso.",
    "e_resultado": "X"
  }
}
```

### `GET /multilogin/2.0.0/servicos/faturas/339800707843/dados-pagamento`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "codigo": "007098175908",
  "protocolo": "20260416280227390",
  "usuario": "***",
  "canalSolicitante": "AGC",
  "distribuidora": "COELBA",
  "regiao": "NE",
  "tipoPerfil": "1",
  "byPassActiv": "X",
  "documentoSolicitante": "*******7586",
  "documento": "*******7586"
}
```
response body:
```json
{
  "codBarras": "838700000052215300300078098175908215056066285939",
  "retorno": {
    "e_resultado": "X"
  }
}
```

### `GET /multilogin/2.0.0/servicos/faturas/339800707843/pdf`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "codigo": "007098175908",
  "protocolo": "20260416280227390",
  "tipificacao": "1031602",
  "usuario": "***",
  "canalSolicitante": "AGC",
  "motivo": "02",
  "distribuidora": "COELBA",
  "regiao": "NE",
  "tipoPerfil": "1",
  "documento": "*******7586",
  "documentoSolicitante": "*******7586",
  "byPassActiv": ""
}
```
response body:
```json
{
  "fileName": "007098175908",
  "fileSize": "62657 ",
  "fileData": "<base64:83544 chars>",
  "fileExtension": ".pdf",
  "retorno": {
    "tipo": "S",
    "id": "ZATCWS",
    "numero": "066",
    "mensagem": "Documento gerado com sucesso.",
    "e_resultado": "X"
  }
}
```

### `GET /multilogin/2.0.0/servicos/faturas/lista-motivo-segundavia`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "usuario": "***",
  "canalSolicitante": "AGC",
  "distribuidora": "COELBA",
  "regiao": "NE",
  "tipoPerfil": "1",
  "documentoSolicitante": "*******7586",
  "codigo": "007098175908"
}
```
response body:
```json
{
  "motivos": [
    {
      "idMotivo": "02",
      "descricao": "NÃO RECEBI - CLIENTE"
    },
    {
      "idMotivo": "03",
      "descricao": "FATURA DANIFICADA - CLIENTE"
    },
    {
      "idMotivo": "04",
      "descricao": "COMPROVAR RESIDÊNCIA - CLIENTE"
    },
    {
      "idMotivo": "07",
      "descricao": "MUDANÇA NA MODALIDADE DE PAGAMENTO"
    },
    {
      "idMotivo": "10",
      "descricao": "NÃO ESTOU COM FATURA EM MÃOS"
    }
  ],
  "retorno": {
    "e_resultado": "X"
  }
}
```

### `GET /multilogin/2.0.0/servicos/faturas/ucs/faturas`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "codigo": "007098175908",
  "documento": "*******7586",
  "canalSolicitante": "AGC",
  "usuario": "***",
  "protocolo": "20260416280227390",
  "tipificacao": "",
  "byPassActiv": "X",
  "documentoSolicitante": "*******7586",
  "documentoCliente": "*******7586",
  "distribuidora": "COELBA",
  "tipoPerfil": "1"
}
```
response body:
```json
{
  "entregaFaturas": {
    "codigoTipoEntrega": "01",
    "descricaoTipoEntrega": "Entrega no Imóvel",
    "enderecoEntrega": "RUA BAHIA 130  LAPAO CENTRO-LAPAO 44905-000",
    "codigoTipoArrecadacao": "01",
    "descricaoTipoArrecadacao": "Indefinida",
    "dataCertaValida": null,
    "dataVencimento": "06",
    "dataCorte": "0000-00-00"
  },
  "faturas": [
    {
      "tipoFatura": {
        "codigo": "PR",
        "descricao": "Periodico"
      },
      "statusFatura": "A Vencer",
      "dataCompetencia": "2026-04-11",
      "dataEmissao": "2026-04-14",
      "dataPagamento": "0000-00-00",
      "dataVencimento": "2026-05-06",
      "grupoTensao": null,
      "mesReferencia": "2026/04",
      "numeroFatura": "339800707843",
      "origemFatura": "FAT",
      "situacaoComercial": "AB",
      "tipoArrecadacao": "Indefinida",
      "tipoEntrega": "Entrega no Imóvel",
      "tipoLeitura": "OSB",
      "uc": "007098175908",
      "valorEmissao": "521.53",
      "dataInicioPeriodo": "2026-03-13",
      "dataFimPeriodo": "2026-04-11",
      "emitidoFatAgrupadora": null,
      "nroFatAgrupadora": null,
      "vencFatAgrupada": "0000-00-00",
      "valorFatAgrupada": "0.00",
      "tipoDoc": "FA",
      "codigoCm": null,
      "numeroBoletoUnico": null,
      "agrupadorContaMinima": null,
      "valorTotalCMAgrupada": "0.00",
      "aceitaPix": "X"
    },
    {
      "tipoFatura": {
        "codigo": "PR",
        "descricao": "Periodico"
      },
      "statusFatura": "Vencida",
      "dataCompetencia": "2026-03-12",
      "dataEmissao": "2026-03-13",
      "dataPagamento": "0000-00-00",
      "dataVencimento": "2026-04-06",
      "grupoTensao": null,
      "mesReferencia": "2026/03",
      "numeroFatura": "337950698598",
      "origemFatura": "FAT",
      "situacaoComercial": "AB",
      "tipoArrecadacao": "Indefinida",
      "tipoEntrega": "Entrega no Imóvel",
      "tipoLeitura": "OSB",
      "uc": "007098175908",
      "valorEmissao": "460.37",
      "dataInicioPeriodo": "2026-02-20",
      "dataFimPeriodo": "2026-03-12",
      "emitidoFatAgrupadora": null,
      "nroFatAgrupadora": null,
      "vencFatAgrupada": "0000-00-00",
      "valorFatAgrupada": "0.00",
      "tipoDoc": "FA",
      "codigoCm": null,
      "numeroBoletoUnico": null,
      "agrupadorContaMinima": null,
      "valorTotalCMAgrupada": "0.00",
      "aceitaPix": "X"
    }
  ],
  "retorno": {
    "e_resultado": "X"
  }
}
```

### `GET /multilogin/2.0.0/servicos/historicos/ucs/007085489032/consumos`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "canalSolicitante": "AGC",
  "usuario": "***",
  "dataInicioPeriodoCalc": "2021-04-18T00:00:00",
  "protocoloSonda": "20260416280227433",
  "opcaoSSOS": "N",
  "protocolo": "20260416280227433",
  "documentoSolicitante": "*******7586",
  "byPassAtiv": "X",
  "distribuidora": "COELBA",
  "tipoPerfil": "1",
  "codigo": "007085489032"
}
```
response body:
```json
{
  "historicoConsumo": [
    {
      "dataPagamento": "06/02/2025",
      "dataVencimento": "08/01/2025",
      "dataLeitura": "19/12/2024",
      "consumoKw": "    96,00",
      "mesReferencia": "12/2024",
      "numeroLeitura": " 8.877,00",
      "tipoLeitura": "01",
      "dataInicioPeriodoCalc": "14/12/2024",
      "dataFimPeriodoCalc": "19/12/2024",
      "dataProxLeitura": "10/01/2025",
      "valorFatura": "113,12",
      "situacaoFatura": "PG",
      "origem": "FAT",
      "numeroFatura": "334075735546",
      "statusFatura": "PAGA",
      "indicativoCustoDisponibilidade": null,
      "indicativoMedia": null,
      "mensagemDisponibilidadeMedia": null,
      "mediadiaria": "16.00"
    },
    {
      "dataPagamento": "20/12/2024",
      "dataVencimento": "20/12/2024",
      "dataLeitura": "13/12/2024",
      "consumoKw": "   433,00",
      "mesReferencia": "12/2024",
      "numeroLeitura": " 8.877,00",
      "tipoLeitura": "01",
      "dataInicioPeriodoCalc": "13/11/2024",
      "dataFimPeriodoCalc": "13/12/2024",
      "dataProxLeitura": "10/01/2025",
      "valorFatura": "536,56",
      "situacaoFatura": "PG",
      "origem": "FAT",
      "numeroFatura": "333775733213",
      "statusFatura": "PAGA",
      "indicativoCustoDisponibilidade": null,
      "indicativoMedia": null,
      "mensagemDisponibilidadeMedia": null,
      "mediadiaria": "13.97"
    },
    {
      "dataPagamento": "20/12/2024",
      "dataVencimento": "21/11/2024",
      "dataLeitura": "12/11/2024",
      "consumoKw": "   424,00",
      "mesReferencia": "11/2024",
      "numeroLeitura": " 8.444,00",
      "tipoLeitura": "01",
      "dataInicioPeriodoCalc": "12/10/2024",
      "dataFimPeriodoCalc": "12/11/2024",
      "dataProxLeitura": "13/12/2024",
      "valorFatura": "582,88",
      "situacaoFatura": "PG",
      "origem": "FAT",
      "numeroFatura": "336850523976",
      "statusFatura": "PAGA",
      "indicativoCustoDisponibilidade": null,
      "indicativoMedia": null,
      "mensagemDisponibilidadeMedia": null,
      "mediadiaria": "13.25"
    },
    {
      "dataPagamento": "01/11/2024",
      "dataVencimento": "18/10/2024",
      "dataLeitura": "11/10/2024",
      "consumoKw": "   378,00",
      "mesReferencia": "10/2024",
      "numeroLeitura": " 8.020,00",
      "tipoLeitura": "01",
      "dataInicioPeriodoCalc": "14/09/2024",
      "dataFimPeriodoCalc": "11/10/2024",
      "dataProxLeitura": "12/11/2024",
      "valorFatura": "497,00",
      "situacaoFatura": "PG",
      "origem": "FAT",
      "numeroFatura": "338325515562",
      "statusFatura": "PAGA",
      "indicativoCustoDisponibilidade": null,
      "indicativoMedia": null,
      "mensagemDisponibilidadeMedia": null,
      "mediadiaria": "13.50"
    },
    {
      "dataPagamento": "01/11/2024",
      "dataVencimento": "20/09/2024",
      "dataLeitura": "13/09/2024",
      "consumoKw": "   677,00",
      "mesReferencia": "09/2024",
      "numeroLeitura": " 7.642,00",
      "tipoLeitura": "01",
      "dataInicioPeriodoCalc": "31/07/2024",
      "dataFimPeriodoCalc": "13/09/2024",
      "dataProxLeitura": "11/10/2024",
      "valorFatura": "808,03",
      "situacaoFatura": "PG",
      "origem": "FAT",
      "numeroFatura": "331450740612",
      "statusFatura": "PAGA",
      "indicativoCustoDisponibilidade": null,
      "indicativoMedia": null,
      "mensagemDisponibilidadeMedia": null,
      "mediadiaria": "15.04"
    }
  ],
  "mediamensal": "    401.6",
  "e_resultado": "X",
  "retorno": null
}
```

### `GET /multilogin/2.0.0/servicos/historicos/ucs/007098175908/consumos`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "canalSolicitante": "AGC",
  "usuario": "***",
  "dataInicioPeriodoCalc": "2021-04-18T00:00:00",
  "protocoloSonda": "20260416280227390",
  "opcaoSSOS": "N",
  "protocolo": "20260416280227390",
  "documentoSolicitante": "*******7586",
  "byPassAtiv": "X",
  "distribuidora": "COELBA",
  "tipoPerfil": "1",
  "codigo": "007098175908"
}
```
response body:
```json
{
  "historicoConsumo": [
    {
      "dataPagamento": "00/00/0000",
      "dataVencimento": "06/05/2026",
      "dataLeitura": "11/04/2026",
      "consumoKw": "   418,00",
      "mesReferencia": "04/2026",
      "numeroLeitura": "15.065,00",
      "tipoLeitura": "01",
      "dataInicioPeriodoCalc": "13/03/2026",
      "dataFimPeriodoCalc": "11/04/2026",
      "dataProxLeitura": "12/05/2026",
      "valorFatura": "521,53",
      "situacaoFatura": "AB",
      "origem": "FAT",
      "numeroFatura": "339800707843",
      "statusFatura": "AVENCER",
      "indicativoCustoDisponibilidade": null,
      "indicativoMedia": null,
      "mensagemDisponibilidadeMedia": null,
      "mediadiaria": "13.93"
    },
    {
      "dataPagamento": "00/00/0000",
      "dataVencimento": "06/04/2026",
      "dataLeitura": "12/03/2026",
      "consumoKw": "   370,00",
      "mesReferencia": "03/2026",
      "numeroLeitura": "14.647,00",
      "tipoLeitura": "01",
      "dataInicioPeriodoCalc": "20/02/2026",
      "dataFimPeriodoCalc": "12/03/2026",
      "dataProxLeitura": "11/04/2026",
      "valorFatura": "460,37",
      "situacaoFatura": "AB",
      "origem": "FAT",
      "numeroFatura": "337950698598",
      "statusFatura": "VENCIDA",
      "indicativoCustoDisponibilidade": null,
      "indicativoMedia": null,
      "mensagemDisponibilidadeMedia": null,
      "mediadiaria": "17.62"
    }
  ],
  "mediamensal": "      394",
  "e_resultado": "X",
  "retorno": null
}
```

### `GET /multilogin/2.0.0/servicos/imoveis/ucs/007085489032`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "usuario": "***",
  "canalSolicitante": "AGC",
  "distribuidora": "COELBA",
  "protocolo": "123",
  "tipoPerfil": "1",
  "opcaoSSOS": "N"
}
```
response body:
```json
{
  "codigo": "007085489032",
  "instalacao": "0009043917",
  "medidor": null,
  "fase": null,
  "local": {
    "endereco": "RUA BAHIA 130",
    "bairro": "CENTRO-LAPAO",
    "codLocalidade": "100000011936",
    "localidade": "LAPAO",
    "codMunicipio": "100000006161",
    "municipio": "LAPAO",
    "cep": "44905-000",
    "localizacao": {
      "sigla": "URBANO",
      "codigo": null,
      "descricao": null
    },
    "uf": "BA",
    "tipoLogradouro": "RU",
    "nomeLogradouro": "RUA BAHIA",
    "numero": "130",
    "complementoEndereco": null
  },
  "enderecoEntrega": {
    "endereco": "RUA BAHIA 130",
    "bairro": "CENTRO-LAPAO",
    "municipio": "LAPAO",
    "cep": "44905-000",
    "uf": "BA"
  },
  "situacao": {
    "codigo": "DS",
    "descricao": "DESLIGADA",
    "dataSituacaoUC": "20241219",
    "cargaInstalada": "470.0000000"
  },
  "dataLigacao": "20240731",
  "cliente": {
    "codigo": "1014945628",
    "nome": "PAULA PEREIRA FERNANDES",
    "dataAtualizacao": "20240731",
    "documento": {
      "tipo": {
        "codigo": "CPF",
        "descricao": "CADASTRO DE PESSOA FISICA"
      },
      "numero": "03021937586"
    },
    "segundoDocumento": {
      "uf": "BA",
      "orgaoExpedidor": {
        "codigo": "SSP",
        "descricao": "Carteira identidade"
      },
      "tipo": {
        "codigo": "Z_RG",
        "descricao": "Carteira identidade"
      },
      "numero": "1386263591"
    },
    "contato": {
      "email": "paulinha_fernandes7@hotmail.com",
      "celular": {
        "ddd": null,
        "numero": null
      },
      "telefone": {
        "ddd": null,
        "numero": "74-981350825"
      }
    },
    "dataNascimento": "19940626"
  },
  "servicos": {
    "baixaRenda": null,
    "faturaEmail": null,
    "dataCerta": null,
    "debitoAutomatico": null,
    "faturaBraile": null,
    "entregaAlternativa": null,
    "debitosVencidos": null,
    "contaMinima": "X",
    "parcelamentoAbertoUc": null
  },
  "caracteristicas": {
    "grandeCliente": null,
    "irrigacao": null,
    "fotovoltaico": null,
    "vip7": null,
    "espelho": null,
    "microMiniGeracaoD": null,
    "tarifaBranca": null
  },
  "mini_micro": null,
  "utd": null,
  "carga": null,
  "e_resultado": "X"
}
```

### `GET /multilogin/2.0.0/servicos/imoveis/ucs/007098175908`
- chamadas observadas: `1`
- status HTTP: `[200]`
```json
{
  "usuario": "***",
  "canalSolicitante": "AGC",
  "distribuidora": "COELBA",
  "protocolo": "123",
  "tipoPerfil": "1",
  "opcaoSSOS": "N"
}
```
response body:
```json
{
  "codigo": "007098175908",
  "instalacao": "0009043917",
  "medidor": "000000001204300816",
  "fase": "MONOFASE",
  "local": {
    "endereco": "RUA BAHIA 130",
    "bairro": "CENTRO-LAPAO",
    "codLocalidade": "100000011936",
    "localidade": "LAPAO",
    "codMunicipio": "100000006161",
    "municipio": "LAPAO",
    "cep": "44905-000",
    "localizacao": {
      "sigla": "URBANO",
      "codigo": null,
      "descricao": null
    },
    "uf": "BA",
    "tipoLogradouro": "RU",
    "nomeLogradouro": "RUA BAHIA",
    "numero": "130",
    "complementoEndereco": null
  },
  "enderecoEntrega": {
    "endereco": "RUA BAHIA 130",
    "bairro": "CENTRO-LAPAO",
    "municipio": "LAPAO",
    "cep": "44905-000",
    "uf": "BA"
  },
  "situacao": {
    "codigo": "LG",
    "descricao": "LIGADA",
    "dataSituacaoUC": "20260220",
    "cargaInstalada": "470.0000000"
  },
  "dataLigacao": "20260220",
  "cliente": {
    "codigo": "1014945628",
    "nome": "PAULA PEREIRA FERNANDES",
    "dataAtualizacao": "20240731",
    "documento": {
      "tipo": {
        "codigo": "CPF",
        "descricao": "CADASTRO DE PESSOA FISICA"
      },
      "numero": "03021937586"
    },
    "segundoDocumento": {
      "uf": "BA",
      "orgaoExpedidor": {
        "codigo": "SSP",
        "descricao": "Carteira identidade"
      },
      "tipo": {
        "codigo": "Z_RG",
        "descricao": "Carteira identidade"
      },
      "numero": "1386263591"
    },
    "contato": {
      "email": "paulinha_fernandes7@hotmail.com",
      "celular": {
        "ddd": null,
        "numero": null
      },
      "telefone": {
        "ddd": null,
        "numero": "74-981350825"
      }
    },
    "dataNascimento": "19940626"
  },
  "servicos": {
    "baixaRenda": null,
    "faturaEmail": null,
    "dataCerta": "X",
    "debitoAutomatico": null,
    "faturaBraile": null,
    "entregaAlternativa": null,
    "debitosVencidos": null,
    "contaMinima": "X",
    "parcelamentoAbertoUc": null
  },
  "caracteristicas": {
    "grandeCliente": null,
    "irrigacao": null,
    "fotovoltaico": null,
    "vip7": null,
    "espelho": null,
    "microMiniGeracaoD": null,
    "tarifaBranca": null
  },
  "mini_micro": null,
  "utd": null,
  "carga": null,
  "e_resultado": "X"
}
```

### `GET /multilogin/2.0.0/servicos/minha-conta`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "canalSolicitante": "AGC",
  "distribuidora": "COELBA",
  "usuario": "*******7586",
  "tipoPerfil": "1",
  "documentoSolicitante": "*******7586"
}
```
response body:
```json
{
  "nome": "Paula Pereira fernandes",
  "usuarioAcesso": "Paula",
  "email": "paulinha_fernandes7@hotmail.com",
  "celular": "(74) 98135-0825",
  "dtUltimaAtualizacao": "2026-04-15T23:42:19.601-03:00",
  "retorno": {
    "tipo": "WSO2",
    "id": "200",
    "numero": "200_OK",
    "mensagem": "Sucesso!"
  }
}
```

### `GET /multilogin/2.0.0/servicos/minha-conta/minha-conta-legado`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "canalSolicitante": "AGC",
  "usuario": "*******7586",
  "usuarioSap": "WSO2_CONEXAO",
  "usuarioSonda": "WSO2_CONEXAO",
  "distribuidora": "COELBA",
  "tipoPerfil": "1",
  "documentoSolicitante": "*******7586"
}
```
response body:
```json
{
  "nomeTitular": "PAULA PEREIRA FERNANDES",
  "clienteDocumentoSecundario": {
    "tipoDocumentoSecundario": "Z_RG",
    "documentoSecundario": "1386263591"
  },
  "dtNascimento": "1994-06-26",
  "emailCadastro": "paulinha_fernandes7@hotmail.com",
  "telefoneContato": "74-981350825",
  "e_resultado": "X",
  "retorno": {
    "tipo": "S",
    "id": "ZATCWS",
    "numero": "029",
    "mensagem": "Executado com sucesso."
  }
}
```

### `GET /polyfills.b50b1d0ce818639d.js`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"\"use strict\";\n(self[\"webpackChunkneoenergia_web\"] = self[\"webpackChunkneoenergia_web\"] || []).push([[6429],{\n\n/***/ 23443:\n/*!**************************!*\\\n  !*** ./src/polyfills.ts ***!\n  \\**************************/\n/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {\n\n__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var _angular_localize_init__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/localize/init */ 14486);\n/* harmony import */ var zone_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! zone.js */ 88583);\n/* harmony import */ var zone_js__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(zone_js__WEBPACK_IMPORTED_MODULE_1__);\n/***************************************************************************************************\n * Load `$localize` onto the global scope - used if i18n tags appear in Angular templates.\n */\n\n/**\n * This file includes polyfills needed by Angular and is loaded before the app.\n * You can add your own extra polyfills to this file.\n *\n * This file is divided into 2 sections:\n *   1. Browser polyfills. These are applied before loading ZoneJS and are sorted by browsers.\n *   2. Application imports. Files imported after ZoneJS that should be loaded before your main\n *      file.\n *\n * The current setup is for so-called \"evergreen\" browsers; the last versions of browsers that\n * automatically update themselves. This includes Safari >= 10, Chrome >= 55 (including Opera),\n * Edge >= 13 on the desktop, and iOS 10 and Chrome on mobile.\n *\n * Learn more in https://angular.io/guide/browser-support\n */\n\n/***************************************************************************************************\n * BROWSER POLYFILLS\n */\n\n/** IE9, IE10 and IE11 requires all of the following polyfills. **/\n// import 'core-js/es6/symbol';\n// import 'core-js/es6/object';\n// import 'core-js/es6/function';\n// import 'core-js/es6/parse-int';\n// import 'core-js/es6/parse-float';\n// import 'core-js/es6/number';\n// import 'core-js/es6/math';\n// import 'core-js/es6/string';\n// import 'core-js/es6/date';\n// import 'core-js/es6/array';\n// import 'core-js/es6/regexp';\n// import 'core-js/es6/map';\n// import 'core-js/es6/weak-map';\n// import 'core-js/es6/set';\n\n/**\n * If the application will be indexed by Google Search, the following is required.\n * Googlebot uses a renderer based on Chrome 41.\n * https://developers.google.com/search/docs/guides/rendering\n **/\n// import 'core-js/es6/array';\n\n/** IE10 and IE11 requires the following for the Reflect API. */\n// import 'core-js/es6/reflect';\n\n/**\n * By default, zone.js will patch all possible macroTask and DomEvents\n * user can disable parts of macroTask/DomEvents patch by setting following flags\n */\n// (window as any).__Zone_disable_requestAnimationFrame = true; // disable patch requestAnimationFrame\n// (window as any).__Zone_disable_on_property = true; // disable patch onProperty such as onclick\n// (window as any).__zone_symbol__BLACK_LISTED_EVENTS = ['scroll', 'mousemove']; // disable patch specified eventNames\n\n/*\n* in IE/Edge developer tools, the addEventListener will also be wrapped by zone.js\n* with the following flag, it will bypass `zone.js` patch for IE/Edge\n*/\n// (window as any).__Zone_enable_cross_context_check = true;\n\n/***************************************************************************************************\n * Zone JS is required by default for Angular itself.\n */\n\n // Included with Angular CLI.\n\n/***************************************************************************************************\n * APPLICATION IMPORTS\n */\n\n/***/ }),\n\n/***/ 88583:\n/*!***********************************************!*\\\n  !*** ./node_modules/zone.js/fesm2015/zone.js ***!\n  \\***********************************************/\n/***/ (() => {\n\n\n/**\n * @license Angular v14.2.0-next.0\n * (c) 2010-2022 Google LLC. https://angular.io/\n * License: MIT\n */\n\n/**\n * @license\n * Copyright Google LLC All Right"
```

### `GET /protocolo/1.1.0/obterProtocolo`
- chamadas observadas: `2`
- status HTTP: `[200]`
```json
{
  "distribuidora": "COEL",
  "canalSolicitante": "AGC",
  "documento": "*******7586",
  "codCliente": "007098175908",
  "recaptchaAnl": "<redacted>",
  "regiao": "NE"
}
```
response body:
```json
{
  "protocoloSalesforce": 20260416280227390,
  "protocoloSalesforceStr": "20260416280227390",
  "protocoloLegado": 20260416280227390,
  "protocoloLegadoStr": "20260416280227390",
  "retorno": {
    "e_resultado": "X"
  }
}
```

### `GET /runtime.a71dfaecea5c5d9f.js`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"/******/ (() => { // webpackBootstrap\n/******/ \t\"use strict\";\n/******/ \tvar __webpack_modules__ = ({});\n/************************************************************************/\n/******/ \t// The module cache\n/******/ \tvar __webpack_module_cache__ = {};\n/******/ \t\n/******/ \t// The require function\n/******/ \tfunction __webpack_require__(moduleId) {\n/******/ \t\t// Check if module is in cache\n/******/ \t\tvar cachedModule = __webpack_module_cache__[moduleId];\n/******/ \t\tif (cachedModule !== undefined) {\n/******/ \t\t\treturn cachedModule.exports;\n/******/ \t\t}\n/******/ \t\t// Create a new module (and put it into the cache)\n/******/ \t\tvar module = __webpack_module_cache__[moduleId] = {\n/******/ \t\t\tid: moduleId,\n/******/ \t\t\tloaded: false,\n/******/ \t\t\texports: {}\n/******/ \t\t};\n/******/ \t\n/******/ \t\t// Execute the module function\n/******/ \t\t__webpack_modules__[moduleId].call(module.exports, module, module.exports, __webpack_require__);\n/******/ \t\n/******/ \t\t// Flag the module as loaded\n/******/ \t\tmodule.loaded = true;\n/******/ \t\n/******/ \t\t// Return the exports of the module\n/******/ \t\treturn module.exports;\n/******/ \t}\n/******/ \t\n/******/ \t// expose the modules object (__webpack_modules__)\n/******/ \t__webpack_require__.m = __webpack_modules__;\n/******/ \t\n/************************************************************************/\n/******/ \t/* webpack/runtime/chunk loaded */\n/******/ \t(() => {\n/******/ \t\tvar deferred = [];\n/******/ \t\t__webpack_require__.O = (result, chunkIds, fn, priority) => {\n/******/ \t\t\tif(chunkIds) {\n/******/ \t\t\t\tpriority = priority || 0;\n/******/ \t\t\t\tfor(var i = deferred.length; i > 0 && deferred[i - 1][2] > priority; i--) deferred[i] = deferred[i - 1];\n/******/ \t\t\t\tdeferred[i] = [chunkIds, fn, priority];\n/******/ \t\t\t\treturn;\n/******/ \t\t\t}\n/******/ \t\t\tvar notFulfilled = Infinity;\n/******/ \t\t\tfor (var i = 0; i < deferred.length; i++) {\n/******/ \t\t\t\tvar [chunkIds, fn, priority] = deferred[i];\n/******/ \t\t\t\tvar fulfilled = true;\n/******/ \t\t\t\tfor (var j = 0; j < chunkIds.length; j++) {\n/******/ \t\t\t\t\tif ((priority & 1 === 0 || notFulfilled >= priority) && Object.keys(__webpack_require__.O).every((key) => (__webpack_require__.O[key](chunkIds[j])))) {\n/******/ \t\t\t\t\t\tchunkIds.splice(j--, 1);\n/******/ \t\t\t\t\t} else {\n/******/ \t\t\t\t\t\tfulfilled = false;\n/******/ \t\t\t\t\t\tif(priority < notFulfilled) notFulfilled = priority;\n/******/ \t\t\t\t\t}\n/******/ \t\t\t\t}\n/******/ \t\t\t\tif(fulfilled) {\n/******/ \t\t\t\t\tdeferred.splice(i--, 1)\n/******/ \t\t\t\t\tvar r = fn();\n/******/ \t\t\t\t\tif (r !== undefined) result = r;\n/******/ \t\t\t\t}\n/******/ \t\t\t}\n/******/ \t\t\treturn result;\n/******/ \t\t};\n/******/ \t})();\n/******/ \t\n/******/ \t/* webpack/runtime/compat get default export */\n/******/ \t(() => {\n/******/ \t\t// getDefaultExport function for compatibility with non-harmony modules\n/******/ \t\t__webpack_require__.n = (module) => {\n/******/ \t\t\tvar getter = module && module.__esModule ?\n/******/ \t\t\t\t() => (module['default']) :\n/******/ \t\t\t\t() => (module);\n/******/ \t\t\t__webpack_require__.d(getter, { a: getter });\n/******/ \t\t\treturn getter;\n/******/ \t\t};\n/******/ \t})();\n/******/ \t\n/******/ \t/* webpack/runtime/create fake namespace object */\n/******/ \t(() => {\n/******/ \t\tvar getProto = Object.getPrototypeOf ? (obj) => (Object.getPrototypeOf(obj)) : (obj) => (obj.__proto__);\n/******/ \t\tvar leafPrototypes;\n/******/ \t\t// create a fake namespace object\n/******/ \t\t// mode & 1: value is a module id, require it\n/******/ \t\t// mode & 2: merge all properties of value into the ns\n/******/ \t\t// mode & 4: return value when already ns object\n/******/ \t\t// mode & 16: return value when it's Promise-like\n/******/ \t\t// mode & 8|1: behave like require\n/******/ \t\t__webpack_require__.t = function(value, mode) {\n/******/ \t\t\tif(mode & 1) value = this(value);\n/******/ \t\t\tif(mode & 8) return value;\n/******/ \t\t\tif(typeof value === 'object' && value) {\n/******/ \t\t\t\tif((mode & 4) && value.__esModule) return value;\n/******/ \t\t\t\tif((mode & 16) && typeof value.then === 'function') return value;\n/******/ \t\t\t}\n/****"
```

### `GET /styles.b8f530f2278fd96d.css`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
"/*!*********************************************************************************************************************************************************************************************************************************************************************************************************************************!*\\\n  !*** css ./node_modules/@angular-devkit/build-angular/node_modules/css-loader/dist/cjs.js??ruleSet[1].rules[4].rules[0].oneOf[0].use[1]!./node_modules/@angular-devkit/build-angular/node_modules/postcss-loader/dist/cjs.js??ruleSet[1].rules[4].rules[0].oneOf[0].use[2]!./node_modules/bootstrap/dist/css/bootstrap.min.css ***!\n  \\*********************************************************************************************************************************************************************************************************************************************************************************************************************************/\n@charset \"UTF-8\";\n/*!\n * Bootstrap  v5.2.3 (https://getbootstrap.com/)\n * Copyright 2011-2022 The Bootstrap Authors\n * Copyright 2011-2022 Twitter, Inc.\n * Licensed under MIT (https://github.com/twbs/bootstrap/blob/main/LICENSE)\n */\n:root{--bs-blue:#0d6efd;--bs-indigo:#6610f2;--bs-purple:#6f42c1;--bs-pink:#d63384;--bs-red:#dc3545;--bs-orange:#fd7e14;--bs-yellow:#ffc107;--bs-green:#198754;--bs-teal:#20c997;--bs-cyan:#0dcaf0;--bs-black:#000;--bs-white:#fff;--bs-gray:#6c757d;--bs-gray-dark:#343a40;--bs-gray-100:#f8f9fa;--bs-gray-200:#e9ecef;--bs-gray-300:#dee2e6;--bs-gray-400:#ced4da;--bs-gray-500:#adb5bd;--bs-gray-600:#6c757d;--bs-gray-700:#495057;--bs-gray-800:#343a40;--bs-gray-900:#212529;--bs-primary:#0d6efd;--bs-secondary:#6c757d;--bs-success:#198754;--bs-info:#0dcaf0;--bs-warning:#ffc107;--bs-danger:#dc3545;--bs-light:#f8f9fa;--bs-dark:#212529;--bs-primary-rgb:13,110,253;--bs-secondary-rgb:108,117,125;--bs-success-rgb:25,135,84;--bs-info-rgb:13,202,240;--bs-warning-rgb:255,193,7;--bs-danger-rgb:220,53,69;--bs-light-rgb:248,249,250;--bs-dark-rgb:33,37,41;--bs-white-rgb:255,255,255;--bs-black-rgb:0,0,0;--bs-body-color-rgb:33,37,41;--bs-body-bg-rgb:255,255,255;--bs-font-sans-serif:system-ui,-apple-system,\"Segoe UI\",Roboto,\"Helvetica Neue\",\"Noto Sans\",\"Liberation Sans\",Arial,sans-serif,\"Apple Color Emoji\",\"Segoe UI Emoji\",\"Segoe UI Symbol\",\"Noto Color Emoji\";--bs-font-monospace:SFMono-Regular,Menlo,Monaco,Consolas,\"Liberation Mono\",\"Courier New\",monospace;--bs-gradient:linear-gradient(180deg, rgba(255, 255, 255, 0.15), rgba(255, 255, 255, 0));--bs-body-font-family:var(--bs-font-sans-serif);--bs-body-font-size:1rem;--bs-body-font-weight:400;--bs-body-line-height:1.5;--bs-body-color:#212529;--bs-body-bg:#fff;--bs-border-width:1px;--bs-border-style:solid;--bs-border-color:#dee2e6;--bs-border-color-translucent:rgba(0, 0, 0, 0.175);--bs-border-radius:0.375rem;--bs-border-radius-sm:0.25rem;--bs-border-radius-lg:0.5rem;--bs-border-radius-xl:1rem;--bs-border-radius-2xl:2rem;--bs-border-radius-pill:50rem;--bs-link-color:#0d6efd;--bs-link-hover-color:#0a58ca;--bs-code-color:#d63384;--bs-highlight-bg:#fff3cd}\n*,::after,::before{box-sizing:border-box}\n@media (prefers-reduced-motion:no-preference){:root{scroll-behavior:smooth}}\nbody{margin:0;font-family:var(--bs-body-font-family);font-size:var(--bs-body-font-size);font-weight:var(--bs-body-font-weight);line-height:var(--bs-body-line-height);color:var(--bs-body-color);text-align:var(--bs-body-text-align);background-color:var(--bs-body-bg);-webkit-text-size-adjust:100%;-webkit-tap-highlight-color:transparent}\nhr{margin:1rem 0;color:inherit;border:0;border-top:1px solid;opacity:.25}\n.h1,.h2,.h3,.h4,.h5,.h6,h1,h2,h3,h4,h5,h6{margin-top:0;margin-bottom:.5rem;font-weight:500;line-height:1.2}\n.h1,h1{font-size:calc(1.375rem + 1.5vw)}\n@media (min-width:1200px){.h1,h1{font-size:2.5rem}}\n.h2,h2{font-size:calc(1.325rem + .9vw)}\n@media (min-width:1200px){.h2,h2{font-size:2rem}}\n.h3,h3{font-size:calc(1.3rem + .6vw)}\n@media (min-width:1200px){."
```

### `GET /undefined`
- chamadas observadas: `1`
- status HTTP: `[200]`
response body:
```json
""
```

### `POST /Z21f/aWAX/7m/5GU3/k93g/3cw3mGrL4mGfJcG71L/Ay4KP0QVAw/IEB5E/QJ1ch4B`
- chamadas observadas: `5`
- status HTTP: `[201]`
request body:
```json
{
  "sensor_data": "3;0;1;0;3486520;FXsoHyqnqJic5h6zRSAPplK96TpZ5bKa147ei2CoUuo=;122,0,0,1,2,0;40ubOPpwN7=9A1|%lfN:[(42h;ypt\"&+S\"/uk#R]3S{S.3S|4\"aN0\"04j05 \".Kd|\"mXMuec\"%WlG\"m~o!eV3Kli\"r!q\"k#7\"Uhre\"STWPw\"SiN\"62!4c,PM\"hzq\"l9Lf0JU_%{+\"x[V\"~{6r(w.U<Q$cLco}94.[ 4QRD)FXp(G?pq6z552)\"M\":(A\"?L-S,PYy$\"HvC\"<7>^f}.N\"C*+\"3\"\")\"XgL\"+\"\"~\"mj5\"z\"nn2In0D]\"J?j\"{g5\"E\"O&h!\"o<x\"Lz`\"E\"I[\"]\"{;{\" \"ND\"se&\"&lD\"n\"\"3\"is\"Z){\"br[Y\"F=~LOa(B&]N~j720sb\".b{\"ht<G iB[y|+eELD}.\"HE|b\"X{G]T?u\"t8Z\"+Yx]o\"0)\"(05dzfs\"KG/!>)/@)p)5O\"J ]\"br\"}\"yW{*\"$\"AF9\"E+[q2}RH]n7N?n\"pP\"/7L\"ho,\"yqV6FU@0\"4I\"0\"Bxk!pC34,Q22YrIq/cl)|+)2(= 0{D8&dM$SAy#O2#k@pj]*H?gw%;CJ>*d2]/(4 [@+y|p 7vcwCZl/`D59<VSa2?+#Z%HLx?-mf\":0+\"?[3\"r\"_WJXUr#6A@;1c\"F\"]vz\")q[^P\"&W+\"V%GSE\"Lq}A\"|7]\"Xoq\"L9l}g\"jf\"}\"<>Vhfaih1`Wb:NJ01J}2J-6&3.\"3\"dhF\"f\"&u}xD{\".m3iq\"5U4\"lPO\"0\"18U\"<\"=[}\"y\"\";\"a(V\"f\"Ab~@Y?,4-6Ji;({&|b{AN8HX#rzwkdS>qGVX?gyYw|GMu}6\"(~?\"M9O\"`oJsF\"~;Y\"s\"^\"%z@x-T?^=\"\"oNp\"Br1\"I\"\"5\"z9U\"h1Z4&QsCT);Tmn\"_|T\"!\"^\"a\"V4l\"@\"NZ\"G]F\"_}H~\"CRi[P(Y\"Tw9\"|\"\"*\"(Em\"Hn8UC\"i6[\"!RV5$\"H,9N\"i=go?\"X?\"ahhV{\"xDj\".N\"Y/P]\"L&&Hc\"eJv\"Z1bIj\"tyo\"1\"01SKH\",h`\"C@MNV\"(iX\"%rZ\"FNb37Pb\"\":\"=iV\"{%+!e\"-cF>\"W\"\"!KR\"<!\"^uqX#\"</LD\"x\"W90Q:&q7r,0kbe/HEeZ@@,r;y/#jV-Ib-nGa,|UQM96H)qvc:gK#7`\"t\".7e\"v\"c7\"~[{\"A?7^\">\"j`MuIiA&1 f$U?/=\"|\"TPa\"+\"|CcwOI2L:kU/bBBEvPGBrS@sRB?*ek5Bc((r(q.kLb0?\"G\" &n\"5;z.i\"5kA\",\"S\"as(\"Wst\"H\"\"c#G\"aeZ\"E\"\"3\"4{@\"Y\"Wh/En\"i~;\":tX\"S\"r@B}QKpa<x(/0[/es4\"K\"~[m\"WkG5`zG31\"38VP:\">|S!8\"jq(\"*cmNY\"%G\"D^!\";5=\"5JX\"eYX\"*S5\"/.s\"3\"\"H\"u)&\"u\"WLLhNY\")RD\"I)_\"P\"X\"5o8\">^4\"D/x[|\"^dyC\"J`,&F/\"Lvrr,K8\"n3~A.\"_L\"U\"\"[%2\"pCN\"g\"\"p\"]]o\" W#^!oq9g<\",^=<\"1=k\"El5\"J?7yA\"BwN \"xP2!Y\"4J[8\"N=e\"MvT\"y\"\"K\"pck\"$\"\"c\"1&r\"Ms&\"ZCA\"OeJyd~dz\"RMK\"b_\"eI/\"D4r\"zWl\"A.:Jg@l\"*M<\""
}
```
response body:
```json
{
  "success": true
}
```

### `POST /akam/13/pixel_32494caf`
- chamadas observadas: `1`
- status HTTP: `[200]`
request body:
```json
"ap=true&bt=%7B%22charging%22%3Atrue%2C%22chargingTime%22%3A0%2C%22dischargingTime%22%3A%22Infinity%22%2C%22level%22%3A1%2C%22onchargingchange%22%3Anull%2C%22onchargingtimechange%22%3Anull%2C%22ondischargingtimechange%22%3Anull%2C%22onlevelchange%22%3Anull%7D&fonts=null&fh=null&timing=%7B%221%22%3A29%2C%222%22%3A1012%2C%22profile%22%3A%7B%22bp%22%3A0%2C%22sr%22%3A1%2C%22dp%22%3A0%2C%22lt%22%3A0%2C%22ps%22%3A0%2C%22cv%22%3A15%2C%22fp%22%3A0%2C%22sp%22%3A3%2C%22br%22%3A0%2C%22ieps%22%3A0%2C%22av%22%3A0%2C%22z1%22%3A8%2C%22jsv%22%3A0%2C%22nav%22%3A0%2C%22nap%22%3A1%2C%22crc%22%3A0%2C%22z2%22%3A9%7D%2C%22main%22%3A378%2C%22compute%22%3A29%2C%22send%22%3A1013%7D&bp=2087755996%2C1953464915%2C591862434%2C325835597%2C1068473606%2C-1382186647%2C-365096851%2C-1979186206%2C-108039040%2C-1906852049&sr=%7B%22inner%22%3A%5B1280%2C900%5D%2C%22outer%22%3A%5B1280%2C900%5D%2C%22screen%22%3A%5B10%2C10%5D%2C%22pageOffset%22%3A%5B0%2C0%5D%2C%22avail%22%3A%5B1280%2C900%5D%2C%22size%22%3A%5B1280%2C900%5D%2C%22client%22%3A%5B1280%2C0%5D%2C%22colorDepth%22%3A24%2C%22pixelDepth%22%3A24%7D&dp=%7B%22XDomainRequest%22%3A0%2C%22createPopup%22%3A0%2C%22removeEventListener%22%3A1%2C%22globalStorage%22%3A0%2C%22openDatabase%22%3A0%2C%22indexedDB%22%3A1%2C%22attachEvent%22%3A0%2C%22ActiveXObject%22%3A0%2C%22dispatchEvent%22%3A1%2C%22addBehavior%22%3A0%2C%22addEventListener%22%3A1%2C%22detachEvent%22%3A0%2C%22fireEvent%22%3A0%2C%22MutationObserver%22%3A1%2C%22HTMLMenuItemElement%22%3A0%2C%22Int8Array%22%3A1%2C%22postMessage%22%3A1%2C%22querySelector%22%3A1%2C%22getElementsByClassName%22%3A1%2C%22images%22%3A1%2C%22compatMode%22%3A%22CSS1Compat%22%2C%22documentMode%22%3A0%2C%22all%22%3A1%2C%22now%22%3A1%2C%22contextMenu%22%3A0%7D&lt=1776388733748-3&ps=true%2Ctrue&cv=e695fc8d72557ec66d8aa490dffcf578688c1825&fp=false&sp=false&br=Chrome&ieps=false&av=false&z=%7B%22a%22%3A843664393%2C%22b%22%3A1%2C%22c%22%3A1%7D&zh=&jsv=1.5&nav=%7B%22userAgent%22%3A%22Mozilla%2F5.0%20(X11%3B%20Linux%20x86_64)%20AppleWebKit%2F537.36%20(KHTML%2C%20like%20Gecko)%20Chrome%2F131.0.0.0%20Safari%2F537.36%22%2C%22appName%22%3A%22Netscape%22%2C%22appCodeName%22%3A%22Mozilla%22%2C%22appVersion%22%3A%225.0%20(X11%3B%20Linux%20x86_64)%20AppleWebKit%2F537.36%20(KHTML%2C%20like%20Gecko)%20Chrome%2F131.0.0.0%20Safari%2F537.36%22%2C%22appMinorVersion%22%3A0%2C%22product%22%3A%22Gecko%22%2C%22productSub%22%3A%2220030107%22%2C%22vendor%22%3A%22Google%20Inc.%22%2C%22vendorSub%22%3A%22%22%2C%22buildID%22%3A0%2C%22platform%22%3A%22Linux%20x86_64%22%2C%22oscpu%22%3A0%2C%22hardwareConcurrency%22%3A16%2C%22language%22%3A%22pt-BR%22%2C%22languages%22%3A%5B%22pt-BR%22%5D%2C%22systemLanguage%22%3A0%2C%22userLanguage%22%3A0%2C%22doNotTrack%22%3Anull%2C%22msDoNotTrack%22%3A0%2C%22cookieEnabled%22%3Atrue%2C%22geolocation%22%3A1%2C%22vibrate%22%3A1%2C%22maxTouchPoints%22%3A0%2C%22webdriver%22%3A0%2C%22plugins%22%3A%5B%22PDF%20Viewer%22%2C%22Chrome%20PDF%20Viewer%22%2C%22Chromium%20PDF%20Viewer%22%2C%22Microsoft%20Edge%20PDF%20Viewer%22%2C%22WebKit%20built-in%20PDF%22%5D%7D&crc=%7B%22window.chrome%22%3A%7B%22app%22%3A%7B%22isInstalled%22%3Afalse%2C%22InstallState%22%3A%7B%22DISABLED%22%3A%22disabled%22%2C%22INSTALLED%22%3A%22installed%22%2C%22NOT_INSTALLED%22%3A%22not_installed%22%7D%2C%22RunningState%22%3A%7B%22CANNOT_RUN%22%3A%22cannot_run%22%2C%22READY_TO_RUN%22%3A%22ready_to_run%22%2C%22RUNNING%22%3A%22running%22%7D%7D%7D%7D&t=1ca27322a4535fa55c559c27156a5496b1bfbe73&u=ce500810ee075130a8eb8bfc2f97bc4e"
```
response body:
```json
""
```

### `POST /areanaologada/2.0.0/autentica`
- chamadas observadas: `1`
- status HTTP: `[200]`
request body:
```json
{
  "usuario": "*******7586",
  "senha": "<redacted>",
  "canalSolicitante": "AGU",
  "recaptcha": "<redacted>"
}
```
response body:
```json
{
  "token": "<redacted>",
  "retorno": {
    "tipo": "WSO2",
    "id": "200",
    "numero": "200_OK",
    "mensagem": "Sucesso"
  }
}
```
