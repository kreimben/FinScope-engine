# FinScope Engine
Back Services for Financial Investment

## 프로젝트 설명
* FinScope에 필요한 각종 데이터들을 전문적으로 수집, 가공 해주는 엔진 서비스 입니다.
* 나에게 어떤 정보가 필요한지 어떤 뉴스와 소식들이 필요한지 알아서 제공해주는 서비스를 만드는 것을 목표로 합니다.
* 금융자산 매도와 매수 그 사이에서 필요한 정보들을 손쉽게 알 수 있는 서비스를 만드는 것을 목표로 합니다.

## 요구사항

1. 데이터 수집 (24시간 작동)
    1. 각종 경제 지표 수집
    1. 펀더멘탈 뉴스 수집
    1. DB에 저장
1. Data 통신
    1. 수집한 데이터를 GraphQL을 통해 클라이언트에 전달

* 경제 뉴스
    * 뉴스 수집 후 Data Cleansing 후 DB에 저장
    * 저장한 데이터는 내부 분석을 위해 활용
    * 이용자에게는 "원문보기"와 같은 기능을 통해 직접 방문하게 만듬
* 경제 지표
    * 경제 지표는 정해진 날 딱 한번 발표 되며 기존의 데이터들은 절대로 변하지 않음
    * 한달에 몇번만, 하지만 실시간으로 빠르게 (발표 후 1분 안에 반영 할 수 있게)

| Status |  Economic Indicator   | Country | Code | Release ID |
|--------|-----------------------|---------|------|------------|
| [x]    | GDP                   | US      | GDP | 53 |
| [x]    | CPI                   | US      | CPIAUCSL | 10 |
| [x]    | Unemployment Rate     | US      | UNRATE | 50 |
| [x]    | M2                    | US      | WM2NS | 21 |
| [x]    | Interest Rate         | US      | DFEDTARU | 101 |
| [x]    | PCE                   | US      | PCEPI | 54 |
| [x]    | Nonfarm Payrolls      | US      | PAYEMS | 50 |
| [x]    | PPI (percent, 전월대비) | US | PPIFIS | 46 |
| [x]    | Initial Claims (주간 실업수당 청구 건수) | US | ICSA | 180 |
| [ ]    | S&P 500               | US      | SP500 |  |
| [ ]    | NASDAQ                | US      | IXIC |  |
| [ ]    | Russell 2000          | US      | RUT |  |
| [ ]    | Dow Jones             | US      | DJIA |  |
| [ ]    | VIX                   | US      | VIXCLS |  |
| [ ]    | 10-Year Treasury Constant Maturity Rate | US | GS10 |  |
| [ ]    | Crude Oil Prices: WTI | US      | DCOILWTICO |  |
| [ ]    | Gold Fixing Price     | US      | GOLDAMGBD228NLBM |  |

* 시장 지수 계산을 위한 기본 데이터 수집

| Status | Index Name |
|--------|------------|
| [ ]    | 공탐지수     |
| [x]    | 버핏지수     |

* 주요 경제 지표 발표 일정

## Tech Stack
* GoLang
* GraphQL
* PostgreSQL
* Supabase
    * Postgres
    * Functions
    * Auth