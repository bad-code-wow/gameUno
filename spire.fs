#version 100

//yea tis like semi-ai coded but i would rather die a horrible death than have to code in glsl without error codes
precision mediump float;

// Input vertex attributes (from vertex shader)
varying vec2 fragTexCoord;
varying vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;
uniform vec2 points[16];

// NOTE: Add your custom variables here

const vec2 size = vec2(800, 450);   // Framebuffer size

const float pi = 3.14159;

vec3 hsvConv(float h){
    float H = h/10.0 - 1.0;
    return vec3(sin(H),sin(H+pi/3.0),sin(H+2.0*pi/3.0));
}

float dist(vec2 p){
    return length(fragTexCoord - p);
}

void main()
{



    // Texel color fetching from texture sampler
    vec4 source = texture2D(texture0, fragTexCoord);
    int best = 0;
    float bestDist = 1e10;
    
    
    for(int i = 0 ; i < 16 ; i++){
    if(dist(points[i])<bestDist - 0.0){
        best = i;
        bestDist = dist(points[i]);
    }
    }

vec3 col = hsvConv(float(best));
    
    gl_FragColor = vec4(col,source.a);

}