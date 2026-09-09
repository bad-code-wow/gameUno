#version 100

precision mediump float;

// Input vertex attributes (from vertex shader)
varying vec2 fragTexCoord;
varying vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;
uniform float time;

// NOTE: Add your custom variables here

const vec2 size = vec2(800, 450);   // Framebuffer size
const float samples = 10.0;          // Pixels per axis; higher = bigger glow, worse performance
const float quality = 1.0;          // Defines size factor: Lower = smaller glow, better quality

float dist(vec2 p){
return pow(pow(p.x,2.0)+pow(p.y,2.0),0.5);
}
float angle(vec2 p,float a){
return fract((atan(p.x,p.y)*5.0 + a)/3.14)*3.14;
}


void main()
{
    vec4 sum = vec4(0);
    vec2 sizeFactor = vec2(1)/size*quality;

    // Texel color fetching from texture sampler
    vec4 source = texture2D(texture0, fragTexCoord);

    float R = dist(fragTexCoord);
    float theta = angle(fragTexCoord,time*2.0+source.r);

    vec3 col = vec3(fract(abs(R*4.0-theta)),fract(abs(R*5.0-theta)),fract(abs(R*1.0-theta)));

    // Calculate final fragment color
    if(source.r == 1 && source.g == 1 && source.b == 1){
    gl_FragColor = vec4(col,1);
    }else{
    gl_FragColor = source;
    }
}