package {{.PackageName}};

{{$hasAp4k := .HasAp4k}}

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;{{if $hasAp4k}}
import io.ap4k.kubernetes.annotation.KubernetesApplication;
import io.ap4k.openshift.annotation.EnableS2iBuild;
import io.ap4k.openshift.annotation.OpenshiftApplication;
{{end}}

@SpringBootApplication{{if $hasAp4k}}
@KubernetesApplication
@EnableS2iBuild
@OpenshiftApplication{{end}}
public class DemoApplication {

    public static void main(String[] args) {
        SpringApplication.run(DemoApplication.class, args);
    }
}
