package model

var LoginTemplate = `<!DOCTYPE html>
<html>

<head>
    <link href='https://fonts.googleapis.com/css?family=IBM Plex Sans Condensed' rel='stylesheet'>
    <style>
        body {
            font-family: 'IBM Plex Sans Condensed';
            font-size: 18px;
            color: #434245;
        }

        a {
            color: inherit;
        }
    </style>
</head>

<body>
    <div style="width: 100%;">
        <table style="width: 446px;margin: 0 auto;">
            <tbody>
                <tr>
                    <td>
                        <table>
                            <tbody>
                                <tr>
                                    <td>
                                        <div>
                                            <div>
                                                <h1
                                                    style="margin-bottom: 35px; font-family:IBM Plex Sans Condensed;font-size: 26px; color: #434245;">
                                                    Bynar</h1>
                                                <h1
                                                    style="margin: 26px 0;font-family:IBM Plex Sans Condensed;font-size: 26px; color: #434245;">
                                                    Confirm your email
                                                    address</h1>
                                                <p style="
       font-size: 20px;
       font-family:IBM Plex Sans Condensed;
       line-height: 28px;
    letter-spacing: -0.2px;
    margin-bottom: 28px;
    word-break: break-word;
">
                                                    Your confirmation code is below and will help you get
                                                    signed in.</p>
                                            </div>
                                            <div style="background: #f5f4f5;
                                                                    width: 100%;
                                                                    margin-left: auto;
                                                                    /* padding-top: 27px; */
                                                                    margin-right: auto;
                                                                    text-align: center;
                                                                    height: 131px;
                                                                    margin-bottom: 30px;">
                                                <p
                                                    style="padding: 43px 0px;
                                                    margin: 0;font-size: 32px;font-family:IBM Plex Sans Condensed;">
                                                    {{.Code}}
                                                </p>
                                            </div>
                                        </div>
                                        <div style="font-size: 16px;font-family:IBM Plex Sans Condensed;">

                                            <p
                                                style="font-size: 16px;font-family:IBM Plex Sans Condensed; color: #434245;margin: 0;">
                                                If you didn't request this email, you can safely ignore it.</p>
                                        </div>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </td>
                </tr>
                <tr>
                    <td>
                        <table>
                            <tbody>
                                <tr>
                                    <td>
                                        <div style="margin-top: 20%;
                                        font-size: 12px;
                                        font-family:IBM Plex Sans Condensed;
                                        color: #696969;
                                        opacity: 80%; ;
                                ">
                                            <div>
                                                <div><a style="font-family:IBM Plex Sans Condensed;font-size: 12px;color: #15c"
                                                        href="https://slackhq.com/" target="_blank"
                                                        data-saferedirecturl="https://www.google.com">Our
                                                        Blog</a>&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;<a
                                                        style="font-family:IBM Plex Sans Condensed;font-size: 12px;color: #15c"
                                                        href="https://google.com" target="_blank"
                                                        data-saferedirecturl="https://www.google.com">Policies</a>&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;<a
                                                        style="font-family:IBM Plex Sans Condensed;font-size: 12px;color: #15c"
                                                        href="https://slack.com/help" target="_blank"
                                                        data-saferedirecturl="https://www.google.com">Help
                                                        Center</a>&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;<a
                                                        style="font-family:IBM Plex Sans Condensed;font-size: 12px;color: #15c"
                                                        class="m_-2755750553702007781footer_link"
                                                        href="https://www.google.com" target="_blank"
                                                        data-saferedirecturl="https://www.google.com">
                                                        Community</a><br /><br />
                                                    <div
                                                        style="font-family:IBM Plex Sans Condensed;font-size: 12px;color: #696969">
                                                        &copy;2023
                                                        Bynar
                                                        415
                                                        Missions Street,
                                                        San Francisco
                                                        </div>
                                                    <br />All rights
                                                    reserved.
                                                </div>
                                            </div>
                                        </div>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</body>

</html>`
