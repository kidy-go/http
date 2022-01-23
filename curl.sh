#!/bin/sh
curl 'http://localhost:9527?from=en&to=zh' \
  -H 'Connection: keep-alive' \
  -H 'Pragma: no-cache' \
  -H 'Cache-Control: no-cache' \
  -H 'sec-ch-ua: " Not;A Brand";v="99", "Google Chrome";v="97", "Chromium";v="97"' \
  -H 'Accept: */*' \
  -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' \
  -H 'X-Requested-With: XMLHttpRequest' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36' \
  -H 'sec-ch-ua-platform: "macOS"' \
  -H 'Origin: https://fanyi.baidu.com' \
  -H 'Sec-Fetch-Site: same-origin' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Referer: https://fanyi.baidu.com/' \
  -H 'Accept-Language: zh-TW,zh;q=0.9,en;q=0.8,zh-CN;q=0.7' \
  -H 'Cookie: PSTM=1618885556; BIDUPSID=4D67FCA335E5F29F057828D9FA4EBEC0; __yjs_duid=1_572b8c5fe34aa41ed1e06f4a6615be6b1622185389076; REALTIME_TRANS_SWITCH=1; FANYI_WORD_SWITCH=1; HISTORY_SWITCH=1; SOUND_SPD_SWITCH=1; SOUND_PREFER_SWITCH=1; BAIDUID=59183672DA19E55F34590083C087B5CD:FG=1; BDSFRCVID_BFESS=4dDOJeC62CJCUqQHYcMvKUL7E3eSMyoTH6aolXlrSugb_hTu_zCYEG0PMx8g0K4M84ljogKK3mOTHR8F_2uxOjjg8UtVJeC6EG0Ptf8g0f5; H_BDCLCKID_SF_BFESS=tJ-toKDhtI83fP36qROq-tuBMfofKRDXKK_shqrpBhcqEIL4htjAynbWyMbJKjciQmtfKRnVXbR8hUbSj4QoDntqDh7wJMckbeoO566zMq5nhMJS257JDMP0-l3OKMJy523iob6vQpnCbhQ3DRoWXPIqbN7P-p5Z5mAqKl0MLPbtbb0xXj_0DTObDGuJJ6KsKjAX3JjV5PK_Hn7zep725M4pbq7H2M-jMHKDbqnJbRTdOU5sjU_KyUPB3Gbn0pcH3mOfhUJb-IOdspcs34tKXTDkQN3T-PRGMIol5b_XatQaDn3oyTbVXp0n0G7ly5jtMgOBBJ0yQ4b4OR5JjxonDh83bG7MJUutfJAjVI_XJID-bnoRq45HMt00qxby26nkBmc9aJ5y-J7nh-cXDq51-nLtM4RPWTb4te3iLfLbQpbZql5O5p7mKnOyX-JZWp5Mte8HKl0MLPboE4nkQxbDetCBKfnMBMPeamOnaU_y3fAKftnOM46JehL3346-35543bRTLnLy5KJYMDFRD5AKej3LjGRabK6aKC5bL6rJabC3qqQoXU6q2bDeQN3EX6oNWI7tBx7cQl4bsDooynj4Dp0vWq54WbbvLT7johRTWqR4eIoODxonDh83hP5O2f5mKJ5XLtjO5hvvhn6O3M7CeMKmDloOW-TB5bbPLUQF5l8-sq0x0bOte-bQXH_EJ50tJJKJ_C-QbRrEDnuz-PvE-PnHMx8X5-RLfbnbKp7F5l8-hlO_Lpoa0fk4-PcqBhQrJ2Ie0xOzMR7xOKQphPvihbD95fJay-vnbGOE5MTN3KJm8tP9bT3vjMrbjMvB2-biWb7L2MbdLDnP_IoG2Mn8M4bb3qOpBtQmJeTxoUJ25DnJhbLGe4bK-Tr-DNDt3D; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; APPGUIDE_10_0_2=1; H_PS_PSSID=35740_35105_31253_35627_35457_34584_35491_35582_35688_26350_35746; BA_HECTOR=2h058g01ag0k0lag4d1gudjeh0r; Hm_lvt_64ecd82404c51e03dc91cb9e8c025574=1642476971,1642477025,1642515931; Hm_lpvt_64ecd82404c51e03dc91cb9e8c025574=1642515931; ab_sr=1.0.1_YjQ5YTQ5NGY3OWUyZDcxNGExZjdhNDZkMDZmYjY4M2EwYTIxYmRhMTMyM2Y0OTllMWY4ZDg2NmIzMGU0MjJmYzAzN2QyOWJlNTJjNWU4NTY4MTNhZmRhYTBlMTM1ZjUx' \
  -H 'x-ticket: eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IlViYVFpdHQ5MUF5WTQwMGMifQ.eyJqdGkiOiJleUowZVhCbElqcHVkV3hzTENKcFpDSTZiblZzYkN3aWRHbHRaU0k2TVRZeU1UTTRPVFkxTVM0MU5qYzBOamg5IiwiaXNzIjoiTHVtZW4iLCJpYXQiOjE2MjEzODk2NTEsImF1ZCI6IjEzIiwiZXhwIjozNzY4ODczMjk4LCJvdl9pZCI6IktlZUd1byJ9.inuRfCtO7OkYmCNZ7fF05nsII_YGqOsYR_94UGy9UUozKT0ukEktUnvirfejP8NWLuMfZetqRecoV3Si6IpGLGCgi3HJ6Uvdn4uDbQRAfWz16ryDG9hXNpJDXTF5bCKXCSVUKk7LF44uqTfLo014mCb533eSAdNTHPCkjCPc-B_MD5mTr2aISB-gqV9A_O8rL78VkWEukkjpUDA6s7D_4eusVWOvZk3Cjoqh9q0Za5pdyPioPJ3ixZIuboevcS6rA3otDF5QBUNfC9qmdUPCN4ZxQVppXb6wZDXt2NTPeCSDlz9un-Hch9OrtJkRnx78M9hzkdJbxWYz5vo_c5yTXQ' \
  -H 'x-dev-debug: true' \
  -D 'from=en&to=zh&query=transaction&transtype=realtime&simple_means_flag=3&sign=155711.426766&token=cc89549df1bdb99a75a2e3d3947374df&domain=common' \
  --compressed
