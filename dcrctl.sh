#!/bin/bash
if [ "$1" == "gettickets" ]; then
    echo '{
    "hashes": [
        "DsQz7MEq1wbgeDZ6WKLd7VLaGYxfT8CkVHX",
        "DsQz7MEq1wbgeDZ6WKLd7VLaGYxfT8CkVHa",
        "DsQz7MEq1wbgeDZ6WKLd7VLaGYxfT8CkVHe",
        "DsYYaFKe3nxWJweGmCaVzPqr2qCa7VePACe",
        "Dsa6UzWBzoCJXE9BtwBiDi8Xd3yunR9yTsf",
        "DsQz7MEq1wbgeDZ6WKLd7VLaGYxfT8CkVfd",
        "DsQz7MEq1wbgeDZ6WKLd7VLaGYxfT8Ckwqe",
        "DsYYaFKe3nxWJweGmCaVzPqr2qCa7Ve43ed",
        "Dsa6UzWBzoCJXE9BtwBiDi8Xd3yunRfkdcf",
        "DsQz7MEq1wbgeDZ6WKLd7VLaGY3eT8CkVHa",
        "DsQz7MEq1wbgeDZ6WKLd7VLaGdcsT8CkVHe",
        "DsYYaFKe3nxWJweGmCaVzPqr3ewd7VePACe",
        "Dsa6UzWBzoCJXE9BtwBiDi8desdfnR9yTsf",
        "DsQz7MEq1wbgeDZ6WKLdsdeaGYxfT8CkVHX",
        "DsQz7MEq1wbgeDZ6WKL45reaGYxfT8CkVHa",
        "DsQz7MEq1wbgeDZ6WK23edfaGYxfT8CkVHe",
        "DsYYaFKe3nxWJweGmCaV4e1r2qCa7VePACe",
        "Dsa6UzWBzoCJXE9Btwdfse34d3yunR9yTsf",
        "DsQz7MEq1wbgeDZ6WK43dfcaGYxfT8CkVfd",
        "DsQz7MEq1wbgeDZ6WKsaw23aGYxfT8Ckwqe",
        "DsYYaFKe3nxWJweGmCgfvb4r2qCa7Ve43ed",
        "Dsa6UzWBzoCJXE9Btwddsw3Xd3yunRfkdcf",
        "DsQz7MEq1wbgeDZ6WKLddsdaGY3eT8CkVHa",
        "DsQz7MEq1wbgeDZ6WK32wxcaGdcsT8CkVHe",
        "DsYYaFKe3nxWJweGmCdsw3gr3ewd7VePACe",
        "Dsa6UzWBzoCJXE9BtwdfgvrdesdfnR9yTsf",
        "DsQz7MEq1wbgeDZ6WKLd7VLaG3wed8CkVHX",
        "DsQz7MEq1wbgeDZ6WKLd7VLaGgfvb8CkVHa",
        "DsQz7MEq1wbgeDZ6WKLd7VLaG1qa38CkVHe",
        "DsYYaFKe3nxWJweGmCaVzPqr2dcs4VePACe",
        "Dsa6UzWBzoCJXE9BtwBiDi8XdhjuyR9yTsf",
        "DsQz7MEq1wbgeDZ6WKLd7VLa4rfvb8CkVfd",
        "DsQz7MEq1wbgeDZ6WKLd7VLaGfg568Ckwqe",
        "DsYYaFKe3nxWJweGmCaVzPqr2ws2hVe43ed",
        "Dsa6UzWBzoCJXE9BtwBiDi8XdkjyxRfkdcf",
        "DsQz7MEq1wbgeDZ6WKLd7VLaGr45d8CkVHa",
        "DsQz7MEq1wbgeDZ6WKLd7VLaG23es8CkVHe",
        "DsYYaFKe3nxWJweGmCaVzPqr3fvb6VePACe",
        "Dsa6UzWBzoCJXE9BtwBiDi8defejhR9yTsf",
        "DsQz7MEq1wbgeDZ6WKLdsdeaG12dv8CkVHX",
        "DsQz7MEq1wbgeDZ6WKL45reaGrt548CkVHa",
        "DsQz7MEq1wbgeDZ6WK23edfaGdv4e8CkVHe",
        "DsYYaFKe3nxWJweGmCaV4e1r2r5vcVePACe",
        "Dsa6UzWBzoCJXE9Btwdfse34d3azxR9yTsf",
        "DsQz7MEq1wbgeDZ6WK43dfcaGdbgr8CkVf0",
        "DsQz7MEq1wbgeDZ6WKsaw23aGfdfr8Ckwqe",
        "DsYYaFKe3nxWJweGmCgfvb4r2de2xVe43ed",
        "Dsa6UzWBzoCJXE9Btwddsw3Xdcfd4Rfkdcf",
        "DsQz7MEq1wbgeDZ6WKLddsdaGnbv58CkVHa",
        "DsQz7MEq1wbgeDZ6WK43dfca3dbgr8CkVfd",
        "DsQz7MEq1wbgeDZ6WKsaw23affdfr8Ckwqe",
        "DsYYaFKe3nxWJweGmCgfvb4rhde2xVe43ed",
        "Dsa6UzWBzoCJXE9Btwddsw3Xvcfd4Rfkdcf",
        "DsQz7MEq1wbgeDZ6WKLddsdasnbv58CkVHa",
        "DsQz7MEq1wbgeDZ6WK32wxcaGsw3c8CkVHe",
        "DsYYaFKe3nxWJweGmCdsw3gr3gfcdVePACe"
    ]
}'
elif [ "$1" == "settspendpolicy" ]; then
echo "settspendpolicy called"
    # echo "Treasury spend policy set for ticket $2: $3"
elif [ "$1" == "getrawmempool" ]; then
    echo '[
        "5c76aa623f8077a075167df35583d652572a62a4e260c1f4085b4edbaa9a5d18"
    ]'
else
    echo "Unrecognized command: $1"
fi
