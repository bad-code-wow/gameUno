
//yea tis like semi-ai coded but i would rather die a horrible death than have to code in glsl without error codes

#version 100

precision mediump float;

varying vec2 fragTexCoord;
varying vec4 fragColor;

uniform sampler2D texture0;
uniform vec4 colDiffuse;

const vec2 size = vec2(800.0, 450.0);
const float quality = 1.0;

void main()
{
    vec2 texel = quality / size;

    vec4 source = texture2D(texture0, fragTexCoord);

    // 5x5 weighted blur: 25 samples instead of 49
    vec4 sum = vec4(0.0);
    float totalWeight = 0.0;

    for (int x = -2; x <= 2; x++)
    {
        for (int y = -2; y <= 2; y++)
        {
            float weight = 1.0 /
                (1.0 + float(x*x + y*y));

            sum += texture2D(
                texture0,
                fragTexCoord + vec2(float(x), float(y)) * texel
            ) * weight;

            totalWeight += weight;
        }
    }

    sum /= totalWeight;

    float vin = abs(fragTexCoord.x - 0.5) * 2.0;

    gl_FragColor = (sum + source * (1.0 - vin)) * colDiffuse;
}
