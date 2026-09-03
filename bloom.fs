#version 100

precision mediump float;

varying vec2 fragTexCoord;
varying vec4 fragColor;

uniform sampler2D texture0;
uniform vec4 colDiffuse;

const vec2 resolution = vec2(1600.0, 900.0);

void main()
{
    vec2 texel = 1.0 / resolution;

    vec4 source = texture2D(texture0, fragTexCoord);
    vec4 blur = vec4(0.0);

    for (int x = -2; x <= 2; x++)
    {
        for (int y = -2; y <= 2; y++)
        {
            vec4 c = texture2D(
                texture0,
                fragTexCoord + vec2(float(x), float(y)) * texel * 3.0
            );

            float brightness = max(c.r, max(c.g, c.b));
            float glow = max(brightness - 0.7, 0.0);

            blur += c * glow;
        }
    }

    blur /= 25.0;

    gl_FragColor = source + blur * 1.5;
}
