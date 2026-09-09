#version 100

precision mediump float;

// Input vertex attributes (from vertex shader)
varying vec2 fragTexCoord;
varying vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;

// NOTE: Add your custom variables here

const vec2 size = vec2(800, 450);   // Framebuffer size
const float samples = 7.0;          // Pixels per axis; higher = bigger glow, worse performance
const float quality = 1.0;          // Defines size factor: Lower = smaller glow, better quality

void main()
{
    vec4 sum = vec4(0);
    vec2 sizeFactor = vec2(1)/size*quality;

    // Texel color fetching from texture sampler
    vec4 source = texture2D(texture0, fragTexCoord);

    const int range = 3;            // should be = (samples - 1)/2;

    for (int x = -range; x <= range; x++)
    {
        for (int y = -range; y <= range; y++)
        {
            sum += texture2D(texture0, fragTexCoord + vec2(x, y)*sizeFactor) * 2.0/(float(x+y)+0.1);
        }
    }


    float vin = abs(fragTexCoord.x-0.5) * 2;

    // Calculate final fragment color
    gl_FragColor = ((sum/(samples*samples)) * (1.0) + source*(1-vin))*colDiffuse;
}