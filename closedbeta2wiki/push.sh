#!/usr/bin/env bash

LOGIN="login"
PASSWORD="password"
SPACE="Space"
ID="1111111"
TEXT="""<table>\
  <tbody>\
    <tr>\
      <th>Заказчик</th>\
      <th>Пул адресов</th>\
    </tr>\
    <tr>\
      <td>\
        #Internal network\
      </td>\
      <td>\
        <p>1.1.1.1</p>\
        <p>2.2.2.2</p>\
      </td>\
    </tr>\
  </tbody>\
</table>"""
NUMBER=16
URL="wikiurl"
\curl -u ${LOGIN}:${PASSWORD} -X PUT -H 'Content-Type: application/json' \
-d"{\"id\":${ID},\"type\":\"page\",\"title\":\"new page\",\"space\":{\"key\":\"${SPACE}\"},\"body\":{\"storage\":{\"value\":\"${TEXT}\",\"representation\":\"storage\"}},\"version\":{\"number\":\"${NUMBER}\"}}" \
${URL} | python -mjson.tool
